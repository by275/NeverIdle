ARG ALPINE_VER=3.23

FROM alpine:${ALPINE_VER} AS alpine
FROM golang:alpine${ALPINE_VER} AS golang
RUN apk add --no-cache git && \
    go build -trimpath -ldflags="-s -w -buildid=" -o /go/bin/noidle github.com/by275/neveridle

FROM alpine
LABEL maintainer="by275"
LABEL org.opencontainers.image.source=https://github.com/by275/NeverIdle

RUN apk add --no-cache tini
COPY --from=golang /go/bin/noidle /usr/local/bin/noidle

# environment settings
ENV LANG=C.UTF-8 \
    PS1="\u@\h:\w\\$ "

ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/noidle"]
