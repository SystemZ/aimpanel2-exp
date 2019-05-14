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

function default {
    up
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}