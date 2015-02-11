FROM golang:1.4
MAINTAINER laiwei laiwei.ustc@gmail.com

ADD . /go/src/github.com/open-falcon/agent
WORKDIR /go/src/github.com/open-falcon/agent
RUN go get
RUN go build

EXPOSE 1988

CMD ./control start && tail -F var/app.log
