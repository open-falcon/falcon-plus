#!/bin/sh

Host=${MYSQL_HOST}
User=${MYSQL_USER}
Password=${MYSQL_PASSWORD}

cd /tmp && \
        git clone --depth=1 https://github.com/open-falcon/falcon-plus && \
        cd /tmp/falcon-plus/ && \
        for x in `ls ./scripts/mysql/db_schema/*.sql`; do
            echo init mysql table $x ...;
            mysql -h${Host} -u${User} -p${Password} < $x;
        done

rm -rf /tmp/falcon-plus/