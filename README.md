Rohr
====

Introduction
------------


Prerequisites
-------------
- Docker (Version 1.13.0 above)
- AWS credentials (including following environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)

Getting Started
---------------
- Download Rohr project:
```sh
git clone git@github.com:concur/rohr.git
```

- Start Rohr project:
```sh
./script/dev up
```

- Stop Rohr project:
You can use `ctl+c`, or use another shell session then run following command
```sh
./script/dev down
```

- Clear Rohr project:
```sh
./script/dev clear
```

- Build eve docker development image:
```sh
./script/dev rebuild
```

- Generate/Update required environment variables:
```sh
./script/dev env
```

Using Eve API Server
--------------------
- Copy following command on Rohr directory
```sh
alias evectl=$PWD/script/evectl 
```

- Change directory to your quoin directory
```sh
cd your/project/infra/dev/quoin
```

- Post/create a new quoin on eve server
```sh
evectl create quoin your_quoin_name
```

- Create a new infrastructure based on your quoin
```sh
evectl create infrastructure your_infrastructure_name using your_quoin_name
```

- Check your infrastructure status
```sh
evectl status your_infrastructure_name
```

- Check your infrastructure state
```sh
evectl state your_infrastructure_name
```

Using evectl
------------

## Prerequisites

- Follow the getting started section above for starting rohr
- In a separate terminal:
  - Make sure you clear your AWS environment variables if they're set
  - Using environment variables or `aws configure`, set AWS access key and secret to your credentials
  - Run `make build`
  - Set the following environment variables:
    - export EVECTL_CA_FILE=$PWD/ca/certs/ca-chain.pem
    - export EVECTL_USERNAME=devop
    - source "$PWD/.devpassword" && export EVECTL_PASSWORD=${devop_pwd}
  - An optional environment variable if you do not want to verify TLS:
    - export EVECTL_TLS_NOVERIFY=true
- Launch a browser and navigate to: http://localhost:8080 and seed the provider table
  - There will be a bootstrap command to do this for us soon

## Provider Authentication

To authenticate against a provider, run the following command. You will then be greeted with questions.

```shell
$ source .envrc
$ evectl authenticate aws
```
