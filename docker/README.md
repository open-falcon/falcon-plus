## Running open-falcon container

`the latest version in docker hub is v0.2.1`

##### 1. Start mysql and init the mysql table before the first running
```
    ## start mysql in container
    docker-compose up -d mysql

    ## init mysql table before the first running, only once needed
    cd /tmp && \
    git clone --depth=1 https://github.com/open-falcon/falcon-plus && \
    cd /tmp/falcon-plus/ && \
    for x in `ls ./scripts/mysql/db_schema/*.sql`; do
        echo init mysql table $x ...;
        docker exec -i falcon-mysql mysql -uroot -ptest123456 < $x;
    done

    rm -rf /tmp/falcon-plus/
```

##### 2. Start redis falcon-plus falcon-dashboard in container
```
    # you can custom your configs in folder 'custom' before starting
    # like: replace the ip '192.168.29.244' with your real server ip
    docker-compose up -d
```

##### 3. Start falcon-plus modules in one container

```
    ## start falcon backend modules, such as graph,api,etc.
    docker exec falcon-plus sh ctrl.sh start \
            graph hbs judge transfer nodata aggregator agent gateway api alarm
    
    ## or you can just start/stop/restart specific module as: 
    docker exec falcon-plus sh ctrl.sh start/stop/restart xxx

    ## check status of backend modules
    docker exec falcon-plus ./open-falcon check
    
    ## or you can check logs at /home/work/open-falcon/logs/ in your host
    ls -l /home/work/open-falcon/logs/
    
```

----

## Building open-falcon images from source code

##### Building falcon-plus

```
    cd /tmp && \
    git clone https://github.com/open-falcon/falcon-plus && \
    cd /tmp/falcon-plus/ && \
    docker build -t falcon-plus:v0.2.1 .
```

##### Building falcon-dashboard
```
    cd /tmp && \
    git clone https://github.com/open-falcon/dashboard  && \
    cd /tmp/dashboard/ && \
    docker build -t falcon-dashboard:v0.2.1 .
```

