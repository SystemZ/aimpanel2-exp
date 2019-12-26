#!/bin/bash
# this script can be run multiple times without damage to OS
[ -z "$1" ] && echo "Provide token as 1st argument" >&2 && exit 1
TOKEN=$1

# settings
AIMPANEL_BINARY_NAME=slave
AIMPANEL_BINARY_DIR=/opt/aimpanel
AIMPANEL_DIR=/opt/aimpanel
# for security we use separate user to run aimpanel and all games
#adduser --system --no-create-home aimpanel
#use add home folder for buildtools and other stuff
adduser --system aimpanel
# create main dirs
mkdir -p $AIMPANEL_DIR/storage
mkdir -p $AIMPANEL_DIR/gs
mkdir -p $AIMPANEL_DIR/trash
# apply permissions
chmod -R 770 $AIMPANEL_DIR/
chown -R aimpanel:root $AIMPANEL_DIR/
wget https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates/latest -O $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
chmod +x $AIMPANEL_BINARY_DIR/$AIMPANEL_BINARY_NAME
# for launching wrapper
ln -s /opt/aimpanel/slave /usr/local/bin/slave
# optitonal java for MC
apt-get install -y openjdk-11-jre-headless
# install service to run aimpanel-slave
echo "
[Unit]
Description=Manage game servers
After=network-online.target

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
Environment="API_URL=https://api-lab.aimpanel.pro"

" > /etc/systemd/system/aimpanel.service
# reload files with services in systemd
systemctl daemon-reload
# start service
systemctl restart aimpanel
# show status to user
systemctl status aimpanel.service
