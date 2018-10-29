#!/bin/sh

DOCKER_DIR=/open-falcon
of_bin=$DOCKER_DIR/open-falcon
DOCKER_HOST_IP=$(route -n | awk '/UG[ \t]/{print $2}')

if [ -z $MYSQL_PORT ]; then
    MYSQL_PORT=$DOCKER_HOST_IP:3306
fi
find $DOCKER_DIR/*/config/*.json -type f -exec sed -i "s/%%MYSQL%%/$MYSQL_PORT/g" {} \;


if [ -z $REDIS_PORT ]; then
    REDIS_PORT=$DOCKER_HOST_IP:6379
fi
find $DOCKER_DIR/*/config/*.json -type f -exec sed -i "s/%%REDIS%%/$REDIS_PORT/g" {} \;

supervisorctl $*
