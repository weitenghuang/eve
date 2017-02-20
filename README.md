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