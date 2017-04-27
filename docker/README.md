## Running falcon-plus container

    docker pull openfalcon/falcon-plus:0.2.0
    docker run -itd -p 8081:8081 openfalcon/falcon-plus:0.2.0 bash /run.sh hbs

## Running falcon-plus container with docker-compose

    docker-compose -f init.yml up -d falcon-plus

## Running mysql and redis container

    docker-compose -f init.yml up -d mysql redis

## Stop and Remove containers

    docker-compose -f init.yml rm -f
