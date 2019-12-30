#!/bin/bash
# docs: https://github.com/adriancooney/Taskfile

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function up {
    docker-compose up -d
    cd master-frontend
    npm run serve
}

function frontend {
    cd master-frontend
    VUE_APP_API_URL=http://localhost:3000 npm run serve
}

# launch this first
function up-slave {
    vagrant up
    vagrant ssh-config > vagrant.slave.ssh.config
}

# should be called by IDE or script after file save in slave codebase
function sync-slave {
    go build -o ./vagrant.slave.bin gitlab.com/systemz/aimpanel2/slave &
    ssh -F vagrant.slave.ssh.config default 'sudo /bin/bash -c "service aimpanel stop"'
    wait
    rsync -e "ssh -F vagrant.slave.ssh.config" --rsync-path="sudo rsync" vagrant.slave.bin default:/opt/aimpanel/slave
    ssh -F vagrant.slave.ssh.config default 'sudo /bin/bash -c "systemctl start aimpanel"'
}

# automatically rebuilds slave and copies binary to VM when slave code changes
# sudo apt-get install -y inotify-tools
function auto-slave {
    inotifywait -e close_write,moved_to,create -r -m ./slave/ |
    while read -r directory events filename; do
      SECONDS=0
      echo -n "syncing..."
      sync-slave
      echo " done in $SECONDS s"
    done
}

function stop {
    docker-compose stop
}

function swagger-gen-dev {
    cd master
    #swagger generate spec -m -o ../swagger.json
    swag init
    mv docs/swagger.json .
    rm -rf docs
    cd ../
}

function swagger-serve {
    swagger serve --flavor=swagger --port=9090 swagger.json
}

function update-todo-issue {
    TODO_CONTENT=$(git grep -n TODO | awk -F':' '{ print "https://gitlab.com/SystemZ/aimpanel2/tree/master/"$1"#L"$2; $1=$2=""; print $0; print "" }' | sed -e 's/^[ \t]*//' | grep -v "Taskfile.sh")
    FIXME_CONTENT=$(git grep -n FIXME | awk -F':' '{ print "https://gitlab.com/SystemZ/aimpanel2/tree/master/"$1"#L"$2; $1=$2=""; print $0; print "" }' | sed -e 's/^[ \t]*//' | grep -v "Taskfile.sh")

    ISSUE_CONTENT="
    ## FIXME
    $FIXME_CONTENT
    ## TODO
    $TODO_CONTENT
    "

    curl -s -o /dev/null -X PUT -H "PRIVATE-TOKEN: $GITLAB_PRIVATE_TOKEN" -F "description=$ISSUE_CONTENT" https://gitlab.com/api/v4/projects/7001381/issues/157
}

function ci-install-rclone-latest {
    wget -q https://downloads.rclone.org/rclone-current-linux-amd64.deb
    dpkg -i rclone-current-linux-amd64.deb
    echo "$RCLONE_CONF" | base64 -d > rclone.conf
}

function ci-slave-upload {
    # we want only production grade stuff available to customers
    if git branch | grep \* | cut -d ' ' -f2 | grep -q master; then
      echo "Master branch detected, deploying files to bucket..."
    else
      echo "Not a master branch, skipping upload to bucket..."
      exit 0
    fi

    ci-install-rclone-latest

    # install.sh is in customer's docs how to install aimpanel
    echo "Uploading install script"
    rclone -v --stats-one-line-date --config rclone.conf copyto install.sh ovh-bucket:aimpanel-updates/install.sh

    # upload as stable link for install.sh
    echo "Uploading binary with \"latest\" name"
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/aimpanel/latest

    # upload as unique link for slave updates, this prevents cache
    echo "Uploading binary with SHA1 name"
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/aimpanel/$CI_COMMIT_SHA

    # we don't need older versions that 1h
    echo "Cleaning binaries older than 1h"
    rclone -v --stats-one-line-date --config rclone.conf delete --min-age 1h ovh-bucket:aimpanel-updates/aimpanel/

    # notify master about new slave version
    # starts all slaves update procedure to be compatible with new master
    # new master will be deployed few minutes after that
    curl -XPOST -H "Token: $AIMPANEL_UPDATE_TOKEN" -H "Content-type: application/json" -d "{\"commit\": \"$CI_COMMIT_SHA\", \"url\": \"https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates/aimpanel/$CI_COMMIT_SHA\"}" 'https://api-lab.aimpanel.pro/v1/version'
}

function ci-redis-compile {
    # based on https://github.com/docker-library/redis/blob/d42494ab2d96070c8d83f37a7542fbbffd999988/5.0/Dockerfile
    apt-get update
    apt-get install -y --no-install-recommends ca-certificates wget gcc libc6-dev make
    wget -O redis.tar.gz "$REDIS_DOWNLOAD_URL"
    echo "$REDIS_DOWNLOAD_SHA *redis.tar.gz" | sha256sum -c -
    mkdir -p /usr/src/redis
    tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
    make -C /usr/src/redis -j "$(nproc)"
    cp /usr/src/redis/src/redis-server redis-server
}

function ci-redis-upload {
    ci-install-rclone-latest
    rclone -v --stats-one-line-date --config rclone.conf copyto redis-server ovh-bucket:aimpanel-updates/redis/$REDIS_VERSION-redis-server-linux-amd64
}

function ci-swagger-gen {
    cd master
    #swagger generate spec -m -o ../swagger.json
    swag init
    mv docs/swagger.json ../
    rm -rf docs
    cd ../
}


function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}