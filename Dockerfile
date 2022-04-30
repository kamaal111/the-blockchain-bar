# syntax=docker/dockerfile:1

FROM ubuntu:20.04
COPY ./bin/tbb /usr/bin/tbb
RUN apt-get update; apt-get install -y curl
CMD /usr/bin/tbb run --datadir=$HOME/.tbb --ip=127.0.0.1 --port=8081 --disable-ssl
