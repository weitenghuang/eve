FROM golang:1.6-alpine
MAINTAINER Concur Platform R&D <platform-engineering@concur.com>

ENV GOPATH="/go"
ENV PATH="$PATH:$GOPATH/bin"
ENV ROHR_PATH="github.com/concur/rohr"

RUN echo "===> Install build dependencies..." \
  && apk --update add --virtual build-dependencies git \
  && apk --update add tar \
  && go get github.com/tools/godep

WORKDIR "$GOPATH/src/$ROHR_PATH"
COPY . "$GOPATH/src/$ROHR_PATH"

RUN echo "===> Install ROHR project dependencies..." \
  && godep restore

RUN echo "===> Run ROHR unit tests..." \
  && go test -v github.com/concur/rohr/...

RUN echo "===> Build ROHR cmd package..." \
  && echo "===> Build eve cmd..." \
  && CGO_ENABLED=0 GOOS=linux go build -i -x -o $GOPATH/bin/eve $ROHR_PATH/cmd/eve

RUN echo "===> Removing build dependencies..." \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

CMD ["tar", "-cvf", "-", "-C", "/go/src/github.com/concur/rohr", "Dockerfile.install", "-C", "/go/bin", "eve"]