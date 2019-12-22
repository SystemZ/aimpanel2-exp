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

function stop {
    docker-compose stop
}

function swagger-gen {
    cd master
    swagger generate spec -o ../swagger.json
    cd ../
}

function swagger-serve {
    swagger serve --flavor=swagger --port=9090 swagger.json
}

function update-todo-issue {
    ISSUE_CONTENT=$(git grep -n  TODO | awk -F':' '{ print "https://gitlab.com/SystemZ/aimpanel2/tree/master/"$1"#L"$2; $1=$2=""; print $0; print "" }' | sed -e 's/^[ \t]*//')
    curl -X PUT -H "PRIVATE-TOKEN: $GITLAB_PRIVATE_TOKEN" -F "description=$ISSUE_CONTENT" https://gitlab.com/api/v4/projects/7001381/issues/157
}

function upload-slave-binary {
    wget https://downloads.rclone.org/rclone-current-linux-amd64.deb
    dpkg -i rclone-current-linux-amd64.deb
    echo "$RCLONE_CONF" | base64 -d > rclone.conf
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/$CI_COMMIT_SHA
    rclone -v --stats-one-line-date --config rclone.conf copyto aimpanel-slave ovh-bucket:aimpanel-updates/latest
    curl -XPOST -H "Token: $AIMPANEL_UPDATE_TOKEN" -H "Content-type: application/json" -d "{\"commit\": \"$CI_COMMIT_SHA\", \"url\": \"https://storage.gra.cloud.ovh.net/v1/AUTH_23b9e96be2fc431d93deedba1b8c87d2/aimpanel-updates/$CI_COMMIT_SHA\"" 'https://api-lab.aimpanel.pro/v1/version'
}

function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}