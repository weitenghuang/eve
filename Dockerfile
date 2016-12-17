FROM golang:1.6-alpine
MAINTAINER Concur Platform R&D <platform-engineering@concur.com>

ENV GOPATH="/go"
ENV PATH="$PATH:$GOPATH/bin:/opt/bin"
ENV ROHR_PATH="github.com/concur/rohr"
ENV TERRAFORM_VERSION=0.7.8
ENV TERRAFORM_SHA256SUM=b3394910c6a1069882f39ad590eead0414d34d5bd73d4d47fa44e66f53454b5a

RUN mkdir -p /opt/bin

RUN apk update

RUN echo "===> Install build dependencies..." \
  && apk --update add --virtual build-dependencies curl jq \
  && apk --update add git \
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

RUN echo "===> Installing Terraform..." \
  && curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip > terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && echo "${TERRAFORM_SHA256SUM}  terraform_${TERRAFORM_VERSION}_linux_amd64.zip" > terraform_${TERRAFORM_VERSION}_SHA256SUMS \
  && sha256sum -cs terraform_${TERRAFORM_VERSION}_SHA256SUMS \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /opt/bin \
  && rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

RUN echo "===> Removing build dependencies..." \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/go/src/github.com/concur/rohr/docker-entrypoint.sh"]

CMD ["tar", "-cvf", "-", "-C", "/go/src/github.com/concur/rohr", "Dockerfile.install", "-C", "/go/bin", "eve",  "-C", "/opt/bin", "terraform"]