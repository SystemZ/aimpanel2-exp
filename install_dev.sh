#!/bin/bash
# this script can be run multiple times without damage to OS
#[ -z "$1" ] && echo "Provide token as 1st argument" >&2 && exit 1
TOKEN="CHANGE_ME"

# debug output
#set -x

# variables
REPO_URL="https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates"
API_URL="http://192.168.122.1:3000"
AIMPANEL_BINARY_NAME=slave
AIMPANEL_BINARY_DIR=/opt/aimpanel
AIMPANEL_DIR=/opt/aimpanel
REDIS_VERSION="5.0.7"
REDIS_BINARY_NAME=redis-server
REDIS_DIR=/opt/aimpanel/redis

# for security we use separate user to run aimpanel and all games
#adduser --system --no-create-home aimpanel
#use add home folder for buildtools and other stuff
adduser --system aimpanel
# create misc dirs
mkdir -p $AIMPANEL_DIR/storage
mkdir -p $AIMPANEL_DIR/gs
mkdir -p $AIMPANEL_DIR/trash
# install redis
mkdir -p $REDIS_DIR
wget "$REPO_URL/redis/$REDIS_VERSION-redis-server-linux-amd64" -O $REDIS_DIR/$REDIS_BINARY_NAME
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
save 60 1
" > $REDIS_DIR/redis.conf

chmod +x $REDIS_DIR/$REDIS_BINARY_NAME
# apply permissions
chmod -R 770 $AIMPANEL_DIR/
chown -R aimpanel:root $AIMPANEL_DIR/
wget $REPO_URL/latest -O $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
chmod +x $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
# for launching wrapper
ln -s /opt/aimpanel/slave /usr/local/bin/slave
# optitonal java for MC
apt-get install -y openjdk-11-jre-headless
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
ExecStart=$REDIS_DIR/$REDIS_BINARY_NAME $REDIS_DIR/redis.conf
WorkingDirectory=/opt/aimpanel/redis/
Restart=always
" > /etc/systemd/system/aimpanel-redis.service
# reload files with services in systemd
systemctl daemon-reload
# start service
systemctl restart aimpanel-redis
# enable autostart
systemctl enable aimpanel-redis
# show status to user
systemctl status --no-pager aimpanel-redis.service

# install service to run aimpanel agent
echo "
[Unit]
Description=Manage game servers
After=network-online.target
Wants=aimpanel-redis.service

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
User=aimpanel
ExecStart=$AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME agent $TOKEN
WorkingDirectory=/opt/aimpanel/
Restart=always
RestartSec=10
Environment="GS_DIR=$AIMPANEL_DIR/gs/"
Environment="STORAGE_DIR=$AIMPANEL_DIR/storage/"
Environment="TRASH_DIR=$AIMPANEL_DIR/trash/"
Environment="HOST_TOKEN=$TOKEN"
Environment="API_URL=$API_URL"
" > /etc/systemd/system/aimpanel.service
# reload files with services in systemd
systemctl daemon-reload
# start service
systemctl restart aimpanel
# enable autostart
systemctl enable aimpanel
# show status to user
systemctl status --no-pager aimpanel.service
