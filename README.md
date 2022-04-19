- [zero](#zero)
  - [设计原则](#%E8%AE%BE%E8%AE%A1%E5%8E%9F%E5%88%99)
  - [跨平台兼容](#%E8%B7%A8%E5%B9%B3%E5%8F%B0%E5%85%BC%E5%AE%B9)
  - [服务注册](#%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C)
  - [服务配置](#%E6%9C%8D%E5%8A%A1%E9%85%8D%E7%BD%AE)
  - [工具安装](#%E5%B7%A5%E5%85%B7%E5%AE%89%E8%A3%85)
  - [安装Protobuff编译环境](#%E5%AE%89%E8%A3%85protobuff%E7%BC%96%E8%AF%91%E7%8E%AF%E5%A2%83)
    - [安装工具](#%E5%AE%89%E8%A3%85%E5%B7%A5%E5%85%B7)
    - [安装命令](#%E5%AE%89%E8%A3%85%E5%91%BD%E4%BB%A4)
  - [快速开始](#%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B)
    - [开启一个nacos-server](#%E5%BC%80%E5%90%AF%E4%B8%80%E4%B8%AAnacos-server)
    - [`greeter.proto`](#greeterproto)
    - [JAVA项目](#java%E9%A1%B9%E7%9B%AE)
      - [演示DEMO](#%E6%BC%94%E7%A4%BAdemo)
      - [生成project 项目](#%E7%94%9F%E6%88%90project-%E9%A1%B9%E7%9B%AE)
      - [生成 Grpc+Protobuf Maven依赖包](#%E7%94%9F%E6%88%90-grpcprotobuf-maven%E4%BE%9D%E8%B5%96%E5%8C%85)
      - [创建一个GRPC服务greeter-service](#%E5%88%9B%E5%BB%BA%E4%B8%80%E4%B8%AAgrpc%E6%9C%8D%E5%8A%A1greeter-service)
      - [创建一个调用服务greeter-api](#%E5%88%9B%E5%BB%BA%E4%B8%80%E4%B8%AA%E8%B0%83%E7%94%A8%E6%9C%8D%E5%8A%A1greeter-api)
    - [GOLANG项目](#golang%E9%A1%B9%E7%9B%AE)
      - [演示Demo](#%E6%BC%94%E7%A4%BAdemo)
      - [创建一个项目](#%E5%88%9B%E5%BB%BA%E4%B8%80%E4%B8%AA%E9%A1%B9%E7%9B%AE)
      - [proto文件管理](#proto%E6%96%87%E4%BB%B6%E7%AE%A1%E7%90%86)
      - [生成greeter-service服务](#%E7%94%9F%E6%88%90greeter-service%E6%9C%8D%E5%8A%A1)
  - [项目结构](#%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84)
    - [JAVA项目结构](#java%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84)
    - [GOLANG项目结构](#golang%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84)
  - [数据校验](#%E6%95%B0%E6%8D%AE%E6%A0%A1%E9%AA%8C)
    - [JAVA数据校验](#java%E6%95%B0%E6%8D%AE%E6%A0%A1%E9%AA%8C)
    - [GOLANG数据校验](#golang%E6%95%B0%E6%8D%AE%E6%A0%A1%E9%AA%8C)
  - [文档生成](#%E6%96%87%E6%A1%A3%E7%94%9F%E6%88%90)
  - [GRPC接口测试](#grpc%E6%8E%A5%E5%8F%A3%E6%B5%8B%E8%AF%95)
  - [熔断降级](#%E7%86%94%E6%96%AD%E9%99%8D%E7%BA%A7)
    - [GOLANG](#golang)
    - [JAVA](#java)
  - [链路追踪](#%E9%93%BE%E8%B7%AF%E8%BF%BD%E8%B8%AA)
  - [授权机制](#%E6%8E%88%E6%9D%83%E6%9C%BA%E5%88%B6)
    - [GOLANG](#golang-1)
    - [JAVA](#java-1)
  - [超时时间](#%E8%B6%85%E6%97%B6%E6%97%B6%E9%97%B4)
    - [GOLANG](#golang-2)
    - [JAVA](#java-2)
  - [其它](#%E5%85%B6%E5%AE%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# zero
本项目是一个微服务开发脚手架，旨在快速构建一个跨平台的微服务架构，帮助开发人员提高开发效率。

## 设计原则

1. 兼容各种语言平台，目前支持 `GOLANG`、`JAVA`。
2. 不重复造轮子，选择社区里比较优秀和成熟的开源框架。
3. 频繁和繁重的工作用代码自动生成提高开发效率。
4. 开源框架的集成做到可扩展性。
5. 统一的约定和规范保障项目的风格规范一致。

## 跨平台兼容

1. 选择 `Grpc`+`Protobuf` 实现RPC服务调用并做到垮语言兼容。
2. 采用`OpenTracing`+`Sleuth`标准实现链路追踪，你可以选择OpenTracing标准下的优秀的链路追踪项目集成。
3. 自研`zctl`帮助工具方便开发者在`Windows`、`Linux`、`Macos`进行代码生成。

## 服务注册

选择`Nacos`作为服务注册中心，框架实现GOLANG版的`Nacos`服务注册和发现。

## 服务配置

选择`Nacos`作为服务配置中心，框架实现GOLANG版的`Nacos`服务配置。

## 工具安装

`go install github.com/SunMaybo/zero/zctl@latest`

![img_1.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img.png)

## 安装Protobuff编译环境

### 安装工具

`protoc-3.20.0`<br>
`protoc-gen-validate-0.6.7`<br>
`protoc-gen-grpc-java-1.45.1`<br>
`protoc-gen-doc_1.5.1`<br>
`protoc-go-inject-tag`<br>
`protoc-gen-go-grpc`<br>

### 安装命令

```
  zctl install --lang golang
```

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img1.png)

**温馨提示**：

1. 你可以通过`--proxy` 指定http_proxy解决下载问题。
2. GOLANG方式安装会帮助安装好JAVA环境所需要的编译插件。

## 快速开始

### 开启一个nacos-server

```
docker run --name nacos-quick -e MODE=standalone -p 8849:8848 -d nacos/nacos-server:2.0.2
```

### `greeter.proto`

```
syntax = "proto3";
package greeter;
option go_package = "/greeter";
import "google/protobuf/timestamp.proto";
option java_package = "com.jewel.meta.asset_platform.proto.greeter_service";
option java_outer_classname = "GreeterServiceProtocol";
import "validate/validate.proto";
service GreeterService {
  rpc SayHelloWord(HelloRequest)returns (HelloReply) {}
  rpc SayStream(stream HelloRequest)returns(stream HelloReply){}
  rpc SayStream1(stream HelloRequest)returns(HelloReply){}
  rpc SayStream2(HelloRequest)returns(stream HelloReply){}
}
message HelloRequest {
  string name = 1;
}
message HelloReply {
  //用户编号
  // @inject_tag: validate:"required,max=32"
  string message = 1 [(validate.rules).string.len = 5];

  google.protobuf.Timestamp time = 2;

}
```

### JAVA项目

#### 演示DEMO

[asset_platform](https://github.com/SunMaybo/asset_platform)

#### 生成project 项目

```
zctl java_project --g com.jewel.meta --a asset_platform
```

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img2.png)

#### 生成 Grpc+Protobuf Maven依赖包

```
zctl java_grpc_package --p ./proto/greeter_service
```

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img3.png)

`注意`： 通过指定Maven 的私服的地址依赖包会上传到私服中，其中：

1. 你可以通过 `--m` 指定Maven的执行路径。
2. 你可以通过`--r` 指定Maven 私服的地址的执行路径。
3. 你可以通过`--v` 指定当前包的版本。
4. 你也可以在`$Home/.zctl/config.yaml `下指定这些配置。

#### 创建一个GRPC服务greeter-service

引入依赖

```
    <dependencies>
        <dependency>
            <groupId>com.jewel.meta.asset_platform</groupId>
            <artifactId>greeter-service-proto</artifactId>
            <version>0.0.2-SNAPSHOT</version>
        </dependency>
    </dependencies>
```

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img4.png)
![img_1.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img5.png)
![img_2.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img6.png)
![img_3.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img7.png)

#### 创建一个调用服务greeter-api

引入依赖

```
    <dependencies>
        <dependency>
            <groupId>com.jewel.meta.asset_platform</groupId>
            <artifactId>greeter-service-proto</artifactId>
            <version>0.0.2-SNAPSHOT</version>
        </dependency>
    </dependencies>
```

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img8.png)

### GOLANG项目

#### 演示Demo

[metdata](https://github.com/SunMaybo/metadata)

#### 创建一个项目

创建一个目录并创建go.mod文件

```
module github.com/SunMaybo/metadata
```

![img_2.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_2.png)

#### proto文件管理

![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img9.png)

#### 生成greeter-service服务

```
/zctl golang_module --m asset_platform --t services
```

![img_1.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_1.png)
![img_3.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_3.png)
![img_4.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_4.png)
**client**
![img_5.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_5.png)
**server**
![img_6.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_6.png)

## 项目结构

### JAVA项目结构

建议一个仓库管理同一个领域的微服务，开发规范按照JAVA开发规范。

```yamml
proto------> 用于存放proto文件
    asset_platform------> 对应模块或者说项目
        - greeter.proto
        - common.proto
greeter_service------> SpringBoot 服务
greeter_api ------> SpringBoot 服务
```

### GOLANG项目结构

建议一个仓库管理多个领域服务。

```yaml
common------> 通用代码目录
proto------> 用于存放proto文件
services
asset_platform------> 对应模块或者说项目
- greeter.proto
- common.proto
apis------> 对外服务的proto文件
asset_api-----> 对应模块或者说项目
-greeter_api.proto
apis------> 对外服务可以是open_api,grpc_web,http，websocket等
asset_api-----> 对应模块或者说项目
- greeter_gateway------>对外服务
- greeter_api
services------>存放对应的微服务
asset_platform------> 对应模块或者说项目
- greeter------> 具体微服务
rpc
etc------>配置文件存放目录
config------> 配置文件mapping层
logic------> RPC服务实现层
svc------>serviceContext 上下文配置
server-----> GRPC 服务层自动生成
main.go-----> 启动入口
tasks----->可以放入一些非接口的任务执行服务
utils----->通用辅助代码目录
tools----->通用工具代码
......

```

## 数据校验

### JAVA数据校验

采用[envovyproxy-validate](https://github.com/envoyproxy/protoc-gen-validate) 进行数据校验，工具自动生成校验代码，你只需要通过配置相应的校验拦截器。

```Java
@Configuration
@RefreshScope
public class GrpcServerAutoConfig {

    @Bean
    public GlobalServerInterceptorConfigurer globalInterceptorConfigurerAdapter() {
        return registry -> {
            registry.add(new ValidatingServerInterceptor(new ReflectiveValidatorIndex()));
        };
    }

}

```

### GOLANG数据校验

采用[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)
通过在pb.go中注入标签进行校验，另外框架封装了[go-playground](github.com/go-playground/validator/v10)进行校验。

## 文档生成

通过工具快速生成接口文档

```
 zctl doc --s ./proto/services/asset_platform
```

## GRPC接口测试

你可以下载[bloomrpc](https://github.com/bloomrpc/bloomrpc) 工具进行接口测试。

## 熔断降级

### GOLANG

目前框架支持go-hystrix进行服务熔断

### JAVA

在SpringBoot体系下你可以灵活选择任何你想使用的工具。

## 链路追踪

采用  `OpenTracing`+`Sleuth`进行链路追踪，目前框架已经支持日志中存放追踪ID信息方便进行问题定位，你也可以通过配置将日志投放到类似`zipkin`链路追踪服务。

## 授权机制

### GOLANG

通过添加拦截器方式做JWT鉴权,并可以过滤掉不需要做鉴权的RPC方法

```
jwtInterceptor := grpc.ChainUnaryInterceptor(
	interceptor.UnaryJWTServerInterceptor("secret", nil),
))
```

GrpcClient调用

```
client,_:=zrpc.NewClient(cfg.Zero).GetGrpcClient("greeter-service", grpc.WithPerRPCCredentials(interceptor.NewTokenAuth("token", false)))
```

### JAVA

可以再调用方和服务实现方通过配置拦截器再metadata中存放authorization：jwt-token，和鉴权Token。

```
@Configuration
@RefreshScope
public class GrpcServerAutoConfig {

    @Bean
    public GlobalServerInterceptorConfigurer globalInterceptorConfigurerAdapter() {
        return registry -> {
            registry.add(new ValidatingServerInterceptor(new ReflectiveValidatorIndex()));
            registry.add(new GlobalGrpcExceptionHandler());
            registry.add(new ServerInterceptor() {
                @Override
                public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(ServerCall<ReqT, RespT> serverCall, Metadata metadata, ServerCallHandler<ReqT, RespT> serverCallHandler) {
                    return null;
                }
            });
        };
    }

}
```

客户端配置授权可以在metadata中写入JWT

```
 greeterServiceBlockingStub.withCallCredentials(new CallCredentials() {
            @Override
            public void applyRequestMetadata(RequestInfo requestInfo, Executor executor, MetadataApplier metadataApplier) {
            }

            @Override
            public void thisUsesUnstableApi() {

            }
        });
```

## 超时时间

### GOLANG

服务端配置整体GRPC调用超时时间

```
zero:
  rpc:
    name: greeter-service
    port: 8088
    is_online: false
    weight: 1
    group_name: format
    cluster_name: default_test
    enable_metrics: true
    timeout: 5000   //RPC服务超时时间配置
    metrics_port: 8843
    metrics_path: /metrics
    metadata:
      active: format
      metrics: "/metrics"
      metrics_port: "8843"
```

客户端指定访问超时时间

```
client,_:=zrpc.NewClient(cfg.Zero).GetGrpcClientWithTimeout(
	"greeter-service", 
	5*time.Second,
	grpc.WithPerRPCCredentials(interceptor.NewTokenAuth("token", false)))
```

### JAVA

不支持

## 其它


