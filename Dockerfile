FROM golang:1.8-alpine 

# LABEL maintainer="laoshancun@foxmail.com"

WORKDIR /usr/local/open-falcon

EXPOSE  8433 6030 5050 8080 8081 6060 5090 1988 9912

CMD ["supervisord","-c","/usr/local/open-falcon/supervisord.conf"]

ENTRYPOINT ["/docker-entrypoint.sh"]

# copies the rest of your code
COPY . /go/src/github.com/open-falcon/falcon-plus/

RUN set -ex \
    && addgroup -S open-falcon && adduser -S -G open-falcon open-falcon \
    # set apk repositories    
    && echo -e "http://mirrors.tuna.tsinghua.edu.cn/alpine/v3.4/main\\nhttp://mirrors.tuna.tsinghua.edu.cn/alpine/v3.4/community" > /etc/apk/repositories \
    # install dependences
    # add bash
    && apk add --update-cache bash supervisor \
    && apk add --virtual .build-deps \
        gcc \
        git \
        make \
        musl-dev \
        py-pip \
    \
    && cd /go/src/github.com/open-falcon/falcon-plus/ \
    && pip install supervisor-stdout \
    # && make misspell \
    && make all \
    && make pack-docker \
    && export VERSION=$(cat VERSION) \
    && mv out/* /usr/local/open-falcon/ \
    && make clean \
    && mv docker-entrypoint.sh / \
    && mv supervisord.conf /usr/local/open-falcon/ \
    && chmod +x /docker-entrypoint.sh \
    && chown -R open-falcon:open-falcon /usr/local/open-falcon/ \
    # cleaning up
    # && rm -rf /go/src/github.com/open-falcon/falcon-plus/ \
    && apk del .build-deps