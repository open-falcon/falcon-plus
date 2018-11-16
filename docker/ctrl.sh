#!/bin/sh

DOCKER_DIR=/open-falcon
of_bin=$DOCKER_DIR/open-falcon
DOCKER_HOST_IP=$(route -n | awk '/UG[ \t]/{print $2}')

#use the correct mysql instance
if [ -z $MYSQL_PORT ]; then
    MYSQL_PORT=$DOCKER_HOST_IP:3306
fi
find $DOCKER_DIR/*/config/*.json -type f -exec sed -i "s/%%MYSQL%%/$MYSQL_PORT/g" {} \;


#use the correct redis instance
if [ -z $REDIS_PORT ]; then
    REDIS_PORT=$DOCKER_HOST_IP:6379
fi
find $DOCKER_DIR/*/config/*.json -type f -exec sed -i "s/%%REDIS%%/$REDIS_PORT/g" {} \;

#use absolute path of metric_list_file in docker
TAB=$'\t'; sed -i "s|.*metric_list_file.*|${TAB}\"metric_list_file\": \"$DOCKER_DIR/api/data/metric\",|g" $DOCKER_DIR/api/config/*.json;

supervisorctl $*
