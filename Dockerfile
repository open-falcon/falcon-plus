FROM golang:1.10-alpine3.7
LABEL maintainer laiwei.ustc@gmail.com
USER root

ENV FALCON_DIR=/open-falcon
ENV PROJ_PATH=${GOPATH}/src/github.com/open-falcon/falcon-plus

RUN mkdir -p $FALCON_DIR && \
    apk add --no-cache ca-certificates bash git g++ perl make supervisor
COPY . ${PROJ_PATH}

WORKDIR ${PROJ_PATH}
RUN make all \
    && make pack4docker \
    && tar -zxf open-falcon-v*.tar.gz -C ${FALCON_DIR} \
    && rm -rf ${PROJ_PATH}
ADD docker/supervisord.conf /etc/supervisord.conf
RUN mkdir -p $FALCON_DIR/logs

EXPOSE 8433 8080
WORKDIR ${FALCON_DIR}

# Start
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
