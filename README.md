# GoIn2.0

### 简介
一款快速开发脚手架框架（基于gin）

## 调试/编译说明

### 调试

1、安装`golang`环境，具体安装过程可以查看官方文档

2、安装 goland 编辑器,并在编辑器设置中配置代理地址为：https://goproxy.io

3、在编辑器中打开该项目并点击右上角运行按钮，即可快速进入调试模式

### 编译

#### 普通编译

1、打开命令行工具并进入到项目根目录，执行 `go build` 命令即可自动编译为当前平台 二进制可执行文件.

#### 交叉编译

##### Mac 下编译， Linux 或者 Windows 下去执行

```
# linux 下去执行
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build ./cmd/flyingstar/main.go
# Windows 下去执行
CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go  build  ./cmd/flyingstar/main.go
```

##### Linux 下编译 ， Mac 或者 Windows 下去执行

```
# Mac  下去执行
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build ./cmd/flyingstar/main.go
# Windows 下执行
CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go build ./cmd/flyingstar/main.go

```

##### Windows 下执行 ， Mac 或 Linux 下去执行

```
# Mac 下执行
SET  CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build ./cmd/flyingstar/main.go

# Linux 去执行
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build ./cmd/flyingstar/main.go

```

### 架构
#### 多入口

我们支持多个启动入口，支持多应用同时被使用，cmd文件夹**goin**为主入口

#### 配置文件

在**config**文件夹下，支持对不同环境下进行对应配置，方便对各个环境的项目部署

#### API文档

```
文档地址gin_swagger

安装swag go get -u github.com/swaggo/swag/cmd/swag

根目录执行 swag init -g ./cmd/goin/main.go -d ./
```

#### 项目业务

我们推荐项目业务编写在**internal**下，在这个文件夹里基础包含：

```
- api 接口主入口
- casbin 权限初始化，采用现在主流的权限验证
- config 配置初始化
- crontab 计时任务
- db 数据库初始化
- globallock 分布式锁
- middleware 中间件
- module 业务模块
 - common 路由初始化模型
 - example 模块
  - model 模型层
  - service 服务层
   - impl 逻辑层
   - pojo 视图模型
   - example.go 服务接口
  - router 路由/控制器
   - application 应用名称
    - v1 版本号（进行版本管理）
  - rdb 模块缓存服务
  - example.go 模块路由注册
 - module.go 路由注册
- notification 消息通知（用于即时信息发布）
- validate 参数验证器
```

#### 消息即时推送

初始化**notification**实现对前端消息即时推送（基于websocket实现），方便于界面实时刷新与数据实时性要求高的通讯

#### utils工具包

包含了一些常用函数，方便使用，**errcode**对项目错误代码集中管理

#### Gorm

我们的模型使用[Gorm2.0](https://gorm.io/zh_CN/docs/index.html)

#### 中间件与权限验证

我们使用redis和[casbin](https://casbin.org/docs/zh-CN/get-started)进行权限验证，基于auth2.0

### 将来要做的

计划 | 备注
---|---
模块自动初始化 | 实现基础业务curd自动创建
组件cli | pkg里的组件按需获取
脚手架cli | 自动创建所需求（api、单功能程序）的架构

组件 | 备注 | 完成度
---|---|---
支付模块 | 包含微信支付、支付宝、普通支付 | 开发中
微信sdk | 集成微信相关接口 | 开发中
企业微信sdk | 集成企业微信相关接口 | 开发中
阿里云sdk | 集成阿里云服务相关接口 | 开发中
抖音sdk | 集成抖音相关接口 | 开发中
飞书sdk | 集成飞书相关接口 | 开发中
七牛sdk | 集成七牛相关接口 | 计划
令牌中心 | token集中管理使用 | 完成
消息即时通知 | 使用websocket技术即时推送消息 | 完成
数据库 | mysql、redis、es... | 补充中（需要补充es）