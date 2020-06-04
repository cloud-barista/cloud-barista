##############################################################
## Stage 1 - Install Go & Set Go env
##############################################################

FROM ubuntu:18.04

RUN apt-get update && yes | apt-get install wget &&\
 wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz &&\
 tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz &&\
 rm go1.13.4.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin


##############################################################
## Stage 2 - Application Set up
##############################################################

# Clone project to docker
ADD . /go/src/github.com/cloud-barista/cb-dragonfly

ENV GOPATH /go/src/github.com/cloud-barista/

WORKDIR $GOPATH/cb-dragonfly

# Use bash
RUN rm /bin/sh && ln -s /bin/bash /bin/sh &&  go mod download && go mod verify

# Run /bin/bash -c "source /app/conf/setup.env"
ENV CBSTORE_ROOT $GOPATH/cb-dragonfly
ENV CBLOG_ROOT $GOPATH/cb-dragonfly
ENV CBMON_PATH $GOPATH/cb-dragonfly

RUN cd $GOPATH/cb-dragonfly/pkg/manager/main;go build -o runMyapp;cp runMyapp $GOPATH/cb-dragonfly

#ENTRYPOINT ["./runMyapp"]
ENV DRAGONFLY_INFLUXDB_URL 127.0.0.1:8086
ENTRYPOINT ["./wait-for-it-wrapper.sh"]

EXPOSE 8094/udp 9090
