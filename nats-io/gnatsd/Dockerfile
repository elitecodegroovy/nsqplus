FROM golang:1.7.6

MAINTAINER Derek Collison <derek@apcera.com>

COPY . /go/src/nats-io/gnatsd
WORKDIR /go/src/nats-io/gnatsd

RUN CGO_ENABLED=0 go install -v -a -tags netgo -installsuffix netgo -ldflags "-s -w -X nats-io/gnatsd/version.GITCOMMIT=`git rev-parse --short HEAD`"

EXPOSE 4222 8222
ENTRYPOINT ["gnatsd"]
CMD ["--help"]
