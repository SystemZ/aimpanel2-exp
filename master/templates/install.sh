#!/bin/bash
# this script can be run multiple times without any problems to OS or panel
# you can use this to repair or upgrade too :)

# variables
TOKEN="{{.Token}}"
REPO_URL="{{.UrlRepo}}"
API_URL="{{.UrlApi}}"
AIMPANEL_BINARY_NAME=slave
AIMPANEL_BINARY_DIR=/opt/aimpanel
AIMPANEL_DIR=/opt/aimpanel
REDIS_VERSION="6.0.6"
REDIS_BINARY_NAME=redis-server
REDIS_DIR=/opt/aimpanel/redis

# init
d_flag=''
verbose='false'
print_usage() {
  printf "Usage: -d (dev)"
}
while getopts ':dv' flag; do
  case "${flag}" in
    d) d_flag='true' ;;
    v) verbose='true' ;;
    *) print_usage
       exit 1 ;;
  esac
done

# debug output
if [ -z "$verbose" ]
then
  set -x
fi

# detect docker
if ! command -v docker &> /dev/null
then
    echo "Docker not found. Please install docker first:"
    echo "wget https://get.docker.com -O- | bash -"
    echo ""
    echo "After docker installation, run it again and at the end run:"
    echo "usermod -aG docker aimpanel"
    exit
fi

# for security we use separate user to run aimpanel and all games
#adduser --system --no-create-home aimpanel
#use add home folder for buildtools and other stuff
adduser --system aimpanel
# create misc dirs
mkdir -p $AIMPANEL_DIR/storage
mkdir -p $AIMPANEL_DIR/gs
mkdir -p $AIMPANEL_DIR/backups
mkdir -p $AIMPANEL_DIR/trash
# install redis
mkdir -p $REDIS_DIR
echo "Downloading redis..."
# turn off redis if already installed
systemctl stop aimpanel-redis
# download redis
wget -q "$REPO_URL/redis/$REDIS_VERSION-redis-server-linux-amd64" -O $REDIS_DIR/$REDIS_BINARY_NAME
# configure redis
echo "
supervised systemd
daemonize no
protected-mode yes
databases 1
port 0
unixsocket $REDIS_DIR/redis.sock
unixsocketperm 700
rdbcompression yes
maxmemory 50M
maxmemory-policy allkeys-lru
tcp-keepalive 30
save ""
save 60 100
" > $REDIS_DIR/redis.conf

chmod +x $REDIS_DIR/$REDIS_BINARY_NAME
# apply permissions
chmod -R 770 $AIMPANEL_DIR/
chown -R aimpanel:root $AIMPANEL_DIR/
echo "Downloading agent..."
wget -q $REPO_URL/aimpanel/latest -O $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
chmod +x $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
# for launching wrapper
ln -s /opt/aimpanel/slave /usr/local/bin/slave
# optitonal java for MC
echo "Updating OS packages..."
apt-get update -q
#echo "Installing systemd integration..."
#apt-get install -q -y libsystemd-dev
echo "Installing OpenJDK 11 JRE..."
apt-get install -q -y openjdk-11-jre-headless

# create service for redis
echo "
[Unit]
Description=DB for aimpanel
After=network-online.target
Before=aimpanel.service

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
User=aimpanel
StandardInput=tty-force
StandardOutput=journal
TTYPath=/dev/tty20
KillMode=process
ExecStart=$REDIS_DIR/$REDIS_BINARY_NAME $REDIS_DIR/redis.conf
WorkingDirectory=/opt/aimpanel/redis/
Restart=always
" > /etc/systemd/system/aimpanel-redis.service

# create service for aimpanel agent
echo "
[Unit]
Description=aimpanel agent
After=network-online.target
Wants=aimpanel-redis.service

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
User=aimpanel
StandardInput=tty-force
TTYPath=/dev/tty20
StandardOutput=journal
KillMode=process
ExecStart=$AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME agent
WorkingDirectory=/opt/aimpanel/
Restart=always
RestartSec=10
ExecStop=$AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME shutdown
Environment="GS_DIR=$AIMPANEL_DIR/gs/"
Environment="STORAGE_DIR=$AIMPANEL_DIR/storage/"
Environment="TRASH_DIR=$AIMPANEL_DIR/trash/"
Environment="HOST_TOKEN=$TOKEN"
Environment="MASTER_URLS=$API_URL"
" > /etc/systemd/system/aimpanel.service

# create service for aimpanel supervisor
echo "
[Unit]
Description=aimpanel supervisor
After=network-online.target
Wants=aimpanel-redis.service

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
User=root
StandardInput=tty-force
TTYPath=/dev/tty20
StandardOutput=journal
KillMode=process
ExecStart=$AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME supervisor
WorkingDirectory=/opt/aimpanel/
Restart=always
RestartSec=10
ExecStop=$AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME shutdown
Environment="GS_DIR=$AIMPANEL_DIR/gs/"
Environment="STORAGE_DIR=$AIMPANEL_DIR/storage/"
Environment="TRASH_DIR=$AIMPANEL_DIR/trash/"
" > /etc/systemd/system/aimpanel-supervisor.service

# reload files with services in systemd
systemctl daemon-reload

# start services
systemctl restart aimpanel-redis
systemctl restart aimpanel-supervisor
systemctl restart aimpanel

# enable autostart
systemctl enable aimpanel-redis
systemctl enable aimpanel-supervisor
systemctl enable aimpanel

# show status to user
#systemctl status --no-pager aimpanel-redis.service
#systemctl status --no-pager aimpanel-supervisor.service
#systemctl status --no-pager aimpanel.service
