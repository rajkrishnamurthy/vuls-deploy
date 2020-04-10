# ---------------------------------------------
FROM golang:alpine as go-cve-dictionary-builder
RUN apk add --no-cache \
        git \
        make \
        gcc \
        musl-dev
ENV REPOSITORY github.com/kotakanbe/go-cve-dictionary
COPY go-cve-dictionary/ $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install
# ---------------------------------------------
FROM golang:alpine as go-exploitdb-builder
RUN apk add --no-cache \
        git \
        make \
        gcc \
        musl-dev
ENV REPOSITORY github.com/vulsio/go-exploitdb
COPY go-exploitdb/ $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install
# ---------------------------------------------
FROM golang:alpine as gost-builder
RUN apk add --no-cache \
        git \
        make \
        gcc \
        musl-dev
ENV REPOSITORY github.com/knqyf263/gost
COPY gost/ $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install
# ---------------------------------------------
FROM golang:alpine as goval-dictionary-builder
RUN apk add --no-cache \
        git \
        make \
        gcc \
        musl-dev
ENV REPOSITORY github.com/kotakanbe/goval-dictionary
COPY goval-dictionary/ $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install
# ---------------------------------------------
FROM golang:alpine as vuls-builder
RUN apk add --no-cache \
        git \
        make \
        gcc \
        musl-dev
ENV REPOSITORY github.com/future-architect/vuls
COPY vuls/ $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install
# ---------------------------------------------

FROM alpine:3.11 as final-builder
LABEL Author="Raj Krishnamurthy" 
ENV LOGDIR /var/log/vuls
ENV WORKDIR /vuls
WORKDIR $WORKDIR
RUN apk add --no-cache ca-certificates \
    && mkdir -p $WORKDIR $LOGDIR

COPY --from=go-cve-dictionary-builder /go/bin/go-cve-dictionary /usr/local/bin/
COPY --from=go-exploitdb-builder /go/bin/go-exploitdb /usr/local/bin/
COPY --from=gost-builder /go/bin/gost /usr/local/bin/
COPY --from=goval-dictionary-builder /go/bin/goval-dictionary /usr/local/bin/
COPY --from=vuls-builder /go/bin/vuls /usr/local/bin/
COPY start.sh start.sh
RUN chmod +x start.sh