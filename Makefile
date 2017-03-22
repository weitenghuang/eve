# Setup environment
PROJECT                 := eve
ORG_PATH                := github.com/concur
REPO_PATH               := ${ORG_PATH}/${PROJECT}
LOCAL_GOPATH            := ${HOME}/.local/share/go/${PROJECT}
GOPATH                  := ${LOCAL_GOPATH}:${GOPATH}
VERSION                 := v0.0.1
ENVIRONMENT             := $${ENVIRONMENT:-DEVELOPMENT}
CMDS                    := $(shell ls ./cmd)

all: .envrc

# Create .envrc populated with environment
.envrc:
	@echo "GOPATH=${GOPATH}" >> .envrc
	@echo "PATH=${LOCAL_GOPATH}/bin:$$PATH" >> .envrc
	@echo "VERSION=${VERSION}" >> .envrc
	@echo "ENVIRONMENT=${ENVIRONMENT}" >> .envrc
	@echo ".envrc file is created at $${PWD}"

# Create .gopath/bin directory
${LOCAL_GOPATH}/bin: .envrc
	@mkdir -p ${LOCAL_GOPATH}/bin

# Create .gopath/src directory
${LOCAL_GOPATH}/src: ${LOCAL_GOPATH}/bin
	@mkdir -p ${LOCAL_GOPATH}/src

# Create .gopath/src/github.com/concur directory
${LOCAL_GOPATH}/src/${ORG_PATH}: ${LOCAL_GOPATH}/src
	@mkdir -p ${LOCAL_GOPATH}/src/${ORG_PATH}

# Create a link from $PWD to .gopath/src/github.com/concur/eve
${LOCAL_GOPATH}/src/${REPO_PATH}: ${LOCAL_GOPATH}/src/${ORG_PATH}
	@ln -s ${PWD} ${LOCAL_GOPATH}/src/${REPO_PATH}

deps: ${LOCAL_GOPATH}/src/${REPO_PATH}
	@glide install

# Build binarys
build: ${LOCAL_GOPATH}/src/${REPO_PATH}
	@for cmd in $(CMDS); do \
		TARGET=$(REPO_PATH)/cmd/$$cmd; \
		go build -o="${LOCAL_GOPATH}/bin/$$cmd" $$TARGET; \
	done

# Run the tests
test: build
	@go test -v $$(glide novendor)

# Cleanup
clean:
	@rm -rf ${LOCAL_GOPATH}

PHONY: all deps build test clean
