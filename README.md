# 🦅 eagle

 [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-eagle/eagle/Go?style=flat-square)](https://github.com/go-eagle/eagle)
 [![codecov](https://codecov.io/gh/go-eagle/eagle/branch/master/graph/badge.svg)](https://codecov.io/gh/go-eagle/eagle)
 [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
 [![godoc](https://godoc.org/github.com/go-eagle/eagle?status.svg)](https://godoc.org/github.com/go-eagle/eagle)
 [![Gitter](https://badges.gitter.im/go-eagle/eagle.svg)](https://gitter.im/go-eagle/eagle?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
 <a href="http://opentracing.io"><img src="https://img.shields.io/badge/OpenTracing-enabled-blue.svg" alt="OpenTracing Badge"></a>
 [![Go Report Card](https://goreportcard.com/badge/github.com/go-eagle/eagle)](https://goreportcard.com/report/github.com/go-eagle/eagle)
 [![gitmoji](https://img.shields.io/badge/gitmoji-%20%F0%9F%98%9C%20%F0%9F%98%8D-FFDD67.svg?style=flat-square)](https://github.com/carloscuesta/gitmoji)
 [![License](https://img.shields.io/github/license/go-eagle/eagle?style=flat-square)](/LICENSE)

一A Go framework suitable for rapid business development, which can quickly build API services or websites.

## official documentation
 - development documentation [https://go-eagle.org/](https://go-eagle.org/)

**Pro Tip:** Basically every directory has `README`，Can make the framework easier to use ^_^

## Design Ideas and Principles

The design ideas and principles used in the framework should satisfy "high cohesion and low coupling" as much as possible, mainly following the following principles
- 1. Single Responsibility Principle
- 2. 基于接口而非实现编程
- 3. Programming based on interfaces rather than implementations
- 4. Multipurpose combination
- 5. Law of Demeter

>Law of Demeter: There should be no dependencies between classes that should not have direct dependencies; between 
> classes with dependencies, try to only rely on the necessary interfaces

## ✨ 技术栈

- Framework routing usage [Gin](https://github.com/gin-gonic/gin) routing
- Middleware usage [Gin](https://github.com/gin-gonic/gin) framework middleware
- Database component [GORM](https://github.com/jinzhu/gorm)
- Documentation usage [Swagger](https://swagger.io/) generate
- Configuration file parsing library [Viper](https://github.com/spf13/viper)
- Use [JWT](https://jwt.io/) Perform identity authentication
- Validator use [validator](https://github.com/go-playground/validator)  It is also the default validator of the Gin framework
- Task scheduling [cron](https://github.com/robfig/cron)
- Package management tool [Go Modules](https://github.com/golang/go/wiki/Modules)
- Test framework [GoConvey](http://goconvey.co/)
- CI/CD [GitHub Actions](https://github.com/actions)
- Use [GolangCI-lint](https://golangci.com/) code detection
- Use make to manage Go project
- Use the shell(admin.sh) script to manage processes
- Multi-environment configuration using YAML files

## 📗 目录结构

```shell
├── Makefile                     # project management file
├── api                          # grpc client and Swagger documentation
├── cmd                          # Scaffolding Directory
├── config                       # Unified storage directory for configuration files
├── docs                         # Framework related documents
├── internal                     # business directory
│   ├── cache                    # Cache based on business encapsulation
│   ├── handler                  # http interface
│   ├── middleware               # custom middleware
│   ├── model                    # database model
│   ├── dao                      # data access layer
│   ├── ecode                    # Business custom error code
│   ├── routers                  # 业务路由
│   ├── server                   # http server and grpc server
│   └── service                  # business logic layer
├── logs                         # directory to store logs
├── main.go                      # project entry file
├── pkg                          # public package
├── test                         # Unit test dependent configuration files, mainly some environment configuration files used by docker
└── scripts                      # Holds scripts for performing various builds, installations, analysis, etc.
```

## 🛠️ Quick start

### method one

The method of Clone project directly, the files are relatively complete

TIPS: MySQL database and Redis need to be installed locally
```bash
# Download and install, you don't need to be GOPATH
git clone https://github.com/go-eagle/eagle

# Go to the download directory
cd eagle

# compile
make build

# run
./scripts/admin.sh start
```

### Method 2

Using scaffolding, only the basic directory is generated, and some public module directories such as pkg are not included

```bash
# download
go get github.com/go-eagle/eagle/cmd/eagle

export GO111MODULE=on
# or in .bashrc or .zshrc
# source .bashrc or source .zshrc

# use
eagle new eagle-demo 
# or 
eagle new github.com/foo/bar
```

## 💻 Common commands

- make helpView help
- make dep Download Go dependencies
- make build Compile the project
- make gen-docs Generate interface documentation
- make test-coverage Generate test coverage
- make lint Check Code Specifications

## 🏂 module

## public module

- Image upload (support local, Qiniu)
- SMS verification code (support Qiniu)

### User module

- register
- Login (email login, mobile phone login)
- Send mobile phone verification code (using Qiniu cloud service)
- Update user information
- follow/unfollow
- Watchlist
- fan list

## 📝 interface documentation

`http://localhost:8080/swagger/index.html`

## development specification

Follow: [Uber Go Language Coding Specification](https://github.com/uber-go/guide/blob/master/style.md)

## 📖 development protocol

- [Configuration instructions](https://github.com/go-eagle/eagle/blob/master/conf)
- [Error code design](https://github.com/go-eagle/eagle/tree/master/pkg/errno)
- [service rules of use](https://github.com/go-eagle/eagle/blob/master/internal/service)
- [repository rules of use](https://github.com/go-eagle/eagle/blob/master/internal/repository)
- [cache Instructions for use](https://github.com/go-eagle/eagle/blob/master/pkg/cache)

## 🚀 deploy

### Deploy separately

After uploading to the server, run the command directly

```bash
./scripts/admin.sh start
```

### Docker deploy

If you have Docker installed, you can start the application with the following command:

```bash
# run
docker-compose up -d

# verify
http://127.0.0.1/health
```

### Supervisord

Compile and generate binaries

```bash
go build -o bin_eagle
```

f the application has multiple machines, it can be compiled on the compilation machine, and then synchronized to 
the corresponding business application server using rsync

> The following can be organized into a script

```bash
export GOROOT=/usr/local/go1.13.8
export GOPATH=/data/build/test/src
export GO111MODULE=on
cd /data/build/test/src/github.com/go-eagle/eagle
/usr/local/go1.13.8/bin/go build -o /data/build/bin/bin_eagle -mod vendor main.go
rsync -av /data/build/bin/ x.x.x.x:/home/go/eagle
supervisorctl restart eagle
```

Here the log directory is set to `/data/log`
If Supervisord is installed, you can add the following to the configuration file (default: `/etc/supervisor/supervisord.conf`):

```ini
[program:eagle]
# environment=
directory=/home/go/eagle
command=/home/go/eagle/bin_eagle
autostart=true
autorestart=true
user=root
stdout_logfile=/data/log/eagle_std.log
startsecs = 2
startretries = 2
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=10
stderr_logfile=/data/log/eagle_err.log
stderr_logfile_maxbytes=10MB
stderr_logfile_backups=10
```

reboot Supervisord

```bash
supervisorctl restart eagle
```

## 📜 CHANGELOG

- [Changelog](https://github.com/go-eagle/eagle/blob/master/CHANGELOG.md)

## 🏘️ who is using

## 💬 Discussion

- Issue: https://github.com/go-eagle/eagle/issues
<img src="https://user-images.githubusercontent.com/3043638/159420999-e00a667d-a5d9-404b-876a-ba0bc94981b9.jpeg" width="200px">

## 📄 License

MIT. See the [LICENSE](LICENSE) file for details.
