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

function swagger-gen {
    cd master
    swagger generate spec -o -m ../swagger.json
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

function ci-slave-installscript-upload {
    rclone -v --stats-one-line-date --config rclone.conf copyto install.sh ovh-bucket:aimpanel-updates/install.sh
}

function ci-slave-binary-upload {
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/$CI_COMMIT_SHA
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/latest
    curl -XPOST -H "Token: $AIMPANEL_UPDATE_TOKEN" -H "Content-type: application/json" -d "{\"commit\": \"$CI_COMMIT_SHA\", \"url\": \"https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates/$CI_COMMIT_SHA\"}" 'https://api-lab.aimpanel.pro/v1/version'
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
    wget https://downloads.rclone.org/rclone-current-linux-amd64.deb
    dpkg -i rclone-current-linux-amd64.deb
    echo "$RCLONE_CONF" | base64 -d > rclone.conf
    rclone -v --stats-one-line-date --config rclone.conf copyto redis-server ovh-bucket:aimpanel-updates/redis/$REDIS_VERSION-redis-server-linux-amd64
}

function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}