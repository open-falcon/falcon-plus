## Owl-light
owl-light is a web UI design for falcon-plus

#### build docker image
```
docker build -t owl-light .
```

#### start web UI
```
docker run -d --name owl-light  -p 80:8080 owl-light -c export declare -x OWL_LIGHT_API_BASE="http://localhost:3000/api/v1"
```

owl-light will use `api` module as backend service.

* After docker container started, [`http://localhost`](http://localhost)
