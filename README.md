[![Build Status](https://travis-ci.com/panicneo/gin_project_starter.svg?branch=master)](https://travis-ci.com/panicneo/gin_project_starter)
[![codecov](https://codecov.io/gh/panicneo/gin_project_starter/branch/master/graph/badge.svg)](https://codecov.io/gh/panicneo/gin_project_starter)
# gin project starter
一个gin项目的初始开发模板

## Feature
 1. web框架使用[gin](github.com/gin-gonic/gin)。
 2. 数据库ORM使用[go-pg](https://github.com/go-pg/pg)，目前模板内只有一个使用PostgreSQL的例子，但是可以通过实现不同的service接口可以轻松采用其他存储源，或者同时使用多种存储，详见storages/services的相关实现。
 3. 配置管理使用[viper](https://github.com/spf13/viper)， 项目模板中使用了toml格式的配置文件，可以轻松改为viper支持的(包括consul/etcd等远程配置方案)任意配置格式。另通过fsevent实现了配置更新时的项目热重载，方便开发调试。
 4. 日志打印统一采用[zerolog](https://github.com/rs/zerolog)，一个高性能的日志库。
 5. 一套实现好的用户系统，包含JWT Token的颁发、用户相关CRUD等API接口。
 6. 基于alpine的两步构建Dockerfile，打包好的镜像在10M左右。

## TODO
 1. [] more test case
 2. [] more feature
 3. [] more docs
