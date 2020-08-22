#!/bin/bash
# docs: https://github.com/adriancooney/Taskfile

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function up {
    # start tmux session with all tools besides running master
    # https://smartosadmin.wordpress.com/2015/03/15/automating-development-workflow-with-tmux/
    SESSION_NAME="exp"
    tmux has-session -t $SESSION_NAME
    if [[ $? != 0 ]]; then
      tmux new-session -s $SESSION_NAME -n dev -d
      # frontend
      #tmux split-window -h -t $SESSION_NAME:0
      tmux send-keys -t $session:0 "./Taskfile.sh frontend" C-m
      # backend deps (mongo + ui) - docker-compose
      tmux split-window -v -t $SESSION_NAME:0.0
      tmux send-keys -t $session:0.1 "docker-compose up" C-m
      # Vagrant VM sync slave
      tmux split-window -v -t $SESSION_NAME:0.1
      tmux send-keys -t $session:0.2 "./Taskfile.sh sync-slave-auto" C-m
      # Vagrant VM
      tmux split-window -v -t $SESSION_NAME:0.2
      tmux send-keys -t $session:0.3 "./Taskfile.sh up-slave; vagrant ssh" C-m
    fi
    tmux a
}

function generate {
    # for "invalid array index" errors
    stringer -type=Id lib/task/task.go
}

function frontend {
    cd master-frontend
    VUE_APP_API_URL=https://127.0.0.1:3000 npm run serve
    # for production use empty "" as API_URL
    #VUE_APP_API_URL="" npm run serve
}

# launch this first
function up-slave {
    vagrant up
    vagrant ssh-config > vagrant.slave.ssh.config
}

function dev-tls-crt {
  openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout key.pem -out crt.pem -config tls.cnf -sha256
}

function update-hosts-slave {
  ssh -F vagrant.slave.ssh.config default "sudo sed -i '/aimpanel.local/d' /etc/hosts"
  ssh -F vagrant.slave.ssh.config default "echo "$(hostname -I | cut -d' ' -f1) aimpanel.local" | sudo tee -a /etc/hosts > /dev/null"
}

# should be called by IDE or script after file save in slave codebase
function sync-slave {
    go build -o ./vagrant.slave.bin gitlab.com/systemz/aimpanel2/slave &
    ssh -F vagrant.slave.ssh.config default 'sudo /bin/bash -c "service aimpanel stop && service aimpanel-supervisor stop"'
    wait
    rsync -e "ssh -F vagrant.slave.ssh.config" --rsync-path="sudo rsync" vagrant.slave.bin default:/opt/aimpanel/slave
    ssh -F vagrant.slave.ssh.config default 'sudo /bin/bash -c "systemctl start aimpanel && systemctl start aimpanel-supervisor"'
}

# automatically rebuilds slave and copies binary to VM when slave code changes
# sudo apt-get install -y inotify-tools
function sync-slave-auto {
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
    TODO_CONTENT=$(git grep -n TODO | grep -v context.TODO | awk -F':' '{ print "https://gitlab.com/SystemZ/aimpanel2/tree/master/"$1"#L"$2; $1=$2=""; print $0; print "" }' | sed -e 's/^[ \t]*//' | grep -v "Taskfile.sh")
    FIXME_CONTENT=$(git grep -n FIXME | awk -F':' '{ print "https://gitlab.com/SystemZ/aimpanel2/tree/master/"$1"#L"$2; $1=$2=""; print $0; print "" }' | sed -e 's/^[ \t]*//' | grep -v "Taskfile.sh")

    ISSUE_CONTENT="
    $FIXME_CONTENT
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
    if [[ "$CI_COMMIT_REF_NAME" == "master" ]]; then
      echo "Master branch detected, deploying files to bucket..."
    else
      echo "Not a master branch, skipping upload to bucket..."
      exit 0
    fi

    ci-install-rclone-latest

    # we don't need this in distributed version with separate instances
    # install.sh is in customer's docs how to install aimpanel
    #echo "Uploading install script"
    #rclone -v --stats-one-line-date --config rclone.conf copyto install.sh ovh-bucket:aimpanel-updates/install.sh

    # upload as stable link for install.sh
    echo "Uploading binary with \"latest\" name"
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/aimpanel/latest

    # upload as unique link for slave updates, this prevents cache
    echo "Uploading binary with SHA1 name"
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/aimpanel/$CI_COMMIT_SHA

    # we don't need older versions than 30d
    echo "Cleaning binaries older than 30d"
    rclone -v --stats-one-line-date --config rclone.conf delete --min-age 30d ovh-bucket:aimpanel-updates/aimpanel/

    # we don't need this in distributed version with separate instances
    # notify master about new slave version
    # starts all slaves update procedure to be compatible with new master
    # new master will be deployed few minutes after that
    #curl -XPOST -H "Token: $AIMPANEL_UPDATE_TOKEN" -H "Content-type: application/json" -d "{\"commit\": \"$CI_COMMIT_SHA\", \"url\": \"https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates/aimpanel/$CI_COMMIT_SHA\"}" 'https://api-lab.aimpanel.pro/v1/version'
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
