from alpine:3.5

env MYSQL_HOST=mysql \
    MYSQL_USER=root \
    MYSQL_PASSWORD=123456

run apk add --no-cache mysql-client git
copy init_mysql_data.sh /

cmd ["/init_mysql_data.sh"]