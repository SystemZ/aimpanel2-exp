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

function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}