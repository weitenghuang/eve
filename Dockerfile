FROM golang:1.8.3-alpine
MAINTAINER Concur Platform R&D <platform-engineering@concur.com>

ENV GOPATH="/go"
ENV PATH="$PATH:$GOPATH/bin:/opt/bin"
ENV EVE_PATH="github.com/concur/eve"
ENV TERRAFORM_VERSION=0.8.7
ENV TERRAFORM_SHA256SUM=7ca424d8d0e06697cc7f492b162223aef525bfbcd69248134a0ce0b529285c8c

RUN mkdir -p /opt/bin && mkdir -p /opt/dist

RUN apk update

RUN echo "===> Install build dependencies..." \
  && apk --update add curl jq git tar zip \
  && git config --global http.https://gopkg.in.followRedirects true \
  && go get github.com/Masterminds/glide

WORKDIR "$GOPATH/src/$EVE_PATH"
COPY . "$GOPATH/src/$EVE_PATH"

RUN echo "===> Install EVE project dependencies..." \
  && glide install

RUN echo "===> Run EVE unit tests..." \
  && go test -v $(glide novendor)

RUN echo "===> Build EVE cmd package..." \
  && mkdir -p $GOPATH/src/$EVE_PATH/artifacts \
  && BUILD_PLATFORM="$(go env GOOS)_$(go env GOARCH)" \
  && echo "===> Build eve cmd..." \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -x -o $GOPATH/src/$EVE_PATH/artifacts/$BUILD_PLATFORM/eve $EVE_PATH/cmd/eve \
  && echo "===> Build evectl cmd..." \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -x -o $GOPATH/src/$EVE_PATH/artifacts/$BUILD_PLATFORM/evectl $EVE_PATH/cmd/evectl \
  && echo "===> Cross Compile evectl" \
  && echo "===> Build evectl cmd for darwin..." \
  && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -i -x -o $GOPATH/src/$EVE_PATH/artifacts/darwin_amd64/evectl $EVE_PATH/cmd/evectl \
  && echo "===> Build evectl cmd for windows..." \
  && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -i -x -o $GOPATH/src/$EVE_PATH/artifacts/windows_amd64/evectl.exe $EVE_PATH/cmd/evectl

RUN echo "===> Packaging evectl..." \
  && for PLATFORM in $(find ./artifacts -mindepth 1 -maxdepth 1 -type d); do \
       OSARCH=$(basename ${PLATFORM}); \
       echo "--> ${OSARCH}"; \
       cd $PLATFORM && zip ../${OSARCH}.zip ./evectl* && cd $OLDPWD; \
     done

RUN echo "===> Preparing distribution..." \
  && for FILENAME in $(find ./artifacts -mindepth 1 -maxdepth 1 -type f); do \
       FILENAME=$(basename $FILENAME); \
       cp ./artifacts/${FILENAME} /opt/dist/evectl_${FILENAME}; \
     done \
  && cd /opt/dist && sha256sum * > ./evectl_SHA256SUMS && cd $OLDPWD \
  && BUILD_PLATFORM="$(go env GOOS)_$(go env GOARCH)" \
  && cp ./artifacts/$BUILD_PLATFORM/eve /opt/bin/eve \
  && cp ./artifacts/$BUILD_PLATFORM/evectl /opt/bin/evectl

RUN echo "===> Installing Terraform..." \
  && curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip > terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && echo "${TERRAFORM_SHA256SUM}  terraform_${TERRAFORM_VERSION}_linux_amd64.zip" > terraform_${TERRAFORM_VERSION}_SHA256SUMS \
  && sha256sum -cs terraform_${TERRAFORM_VERSION}_SHA256SUMS \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /opt/bin \
  && rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

ENTRYPOINT ["/go/src/github.com/concur/eve/docker-entrypoint.sh"]

CMD ["tar", "-cvf", "-", "-C", "/go/src/github.com/concur/eve", "Dockerfile.install", "-C", "/go/src/github.com/concur/eve", "docker-entrypoint.sh", "-C", "/opt/bin", "eve", "-C", "/opt/bin", "evectl", "-C", "/opt/bin", "terraform", "-C", "/opt/dist", "evectl_darwin_amd64.zip", "-C", "/opt/dist", "evectl_windows_amd64.zip", "-C", "/opt/dist", "evectl_linux_amd64.zip", "-C", "/opt/dist", "evectl_SHA256SUMS"]