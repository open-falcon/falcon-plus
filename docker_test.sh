#!/bin/bash

export DB_USER=root
export DB_PASSWORD=test123456
export DB_HOST=127.0.0.1
export DB_PORT=13306
export REDIS_HOST=127.0.0.1
export REDIS_PORT=16379
export API_PORT=18080
export API_HOST=127.0.0.1

docker rm -f falcon-mysql falcon-redis falcon-plus &> /dev/null
if [[ `uname -m` == "aarch64" ]]; then
	docker run --name falcon-mysql -e MYSQL_ROOT_PASSWORD=$DB_PASSWORD -p $DB_PORT:3306 -d mariadb:10.3
else
	docker run --name falcon-mysql -e MYSQL_ROOT_PASSWORD=$DB_PASSWORD -p $DB_PORT:3306 -d mysql:5.7
fi
docker run --name falcon-redis -p $REDIS_PORT:6379 -d redis:4-alpine3.8

echo "waiting mysql start..."
sleep 15
for x in `ls ./scripts/mysql/db_schema/*.sql`; do
    echo "- - -" $x ...
    mysql -h $DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD < $x
done

commit_id=`git rev-parse --short HEAD`
image_tag="falcon-plus:$commit_id"

#build docker image from source code
if [[ `uname -m` == "aarch64" ]]; then
	docker build -t $image_tag -f Dockerfile_arm64 .
else
	docker build -t $image_tag .
fi

## run falcon-plus container
docker run -itd --name falcon-plus \
	 --link=falcon-mysql:db.falcon \
	 --link=falcon-redis:redis.falcon \
	 -p 18433:8433 \
	 -p 18080:8080 \
	 -e MYSQL_PORT=$DB_USER:$DB_PASSWORD@tcp\(db.falcon:3306\) \
	 -e REDIS_PORT=redis.falcon:6379  \
	 $image_tag
sleep 15

## start falcon backend modules, such as graph,api,etc.
docker exec falcon-plus sh ctrl.sh start \
		graph hbs judge transfer nodata aggregator agent gateway api alarm

echo "sleep 15s waiting for falcon-plus process ready..."
sleep 15

make test
