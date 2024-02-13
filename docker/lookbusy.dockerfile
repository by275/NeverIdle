# https://github.com/Fog-Forest/scripts/tree/main/oracle-lookbusy
FROM alpine:3.19 AS alpine

FROM alpine AS lookbusy

WORKDIR /tmp
RUN wget http://www.devin.com/lookbusy/download/lookbusy-1.4.tar.gz && tar -xzvf lookbusy-1.4.tar.gz
RUN apk add build-base libtool
RUN cd ./lookbusy-1.4 && ./configure && make && make install

FROM alpine AS speedtest

WORKDIR /tmp
RUN file="ookla-speedtest-1.2.0-linux-$(apk info --print-arch).tgz" && \
    wget "https://install.speedtest.net/app/cli/$file" && \
    tar xvfz "$file"

FROM alpine AS release
COPY --from=lookbusy /usr/local/bin /usr/local/bin/
COPY --from=speedtest /tmp/speedtest /usr/local/bin/
# SIZE 10.1MB
