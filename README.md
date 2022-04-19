# zero
本项目是一个微服务开发脚手架，旨在快速构建一个跨平台的微服务架构，帮助开发人员提高开发效率。
## 设计原则
1. 兼容各种语言平台，目前支持 `Golang`、`Java`。
2. 不重复造轮子，选择社区里比较优秀和成熟的开源框架。
3. 频发和繁重的工作用代码自动生成提高开发效率。
4.  开源框架集成做到可扩展性。
5.  统一约定和规范保障项目的风格规范一致。

## 跨平台兼容
1. 选择 `Grpc`+`Protobuf` 实现RPC服务调用并做到垮语言兼容。
2. 采用`OpenTracing`+`Sleuth`标准实现链路追踪，你可以选择OpenTracing标准下的优秀的链路追踪项目集成。
3. 自研`zctl`帮助工具方便开发者在`Windows`、`Linux`、`Macos`进行代码生成。

## 服务注册
选择`Nacos`作为服务注册中心，框架实现Golang版的`Nacos`服务注册和发现。
## 服务配置
选择`Nacos`作为服务配置中心，框架实现Golang版的`Nacos`服务配置。

## 工具安装
`go install github.com/SunMaybo/zero/zctl@latest`

![img_1.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img.png)

## 安装Protobuff编译环境
### 安装工具
`protoc-3.20.0`
`protoc-gen-validate-0.6.7`
`protoc-gen-grpc-java-1.45.1`
`protoc-gen-doc_1.5.1`
`protoc-go-inject-tag`
`protoc-gen-go-grpc`

### 安装命令
```
  zctl install --lang golang
```
![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img1.png)

**温馨提示**：
1. 你可以通过`--proxy` 指定http_proxy解决下载问题。
2. golang方式安装会帮助安装好java环境所需要的编译插件。

## 快速开始
### 开启一个Nacos-server
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
## JAVA项目
### 生成project 项目
```
zctl java_project --g com.jewel.meta --a asset_platform
```
![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img2.png)

### 生成 Grpc+Protobuf Maven依赖包
```
zctl java_grpc_package --p ./proto/greeter_service
```
![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img3.png)

`注意`：
通过指定Maven 的私服的地址依赖包会上传到私服中，其中：

1.  你可以通过 `--m` 指定Maven的执行路径。
2.  你可以通过`--r` 指定Maven 私服的地址的执行路径。
3.  你可以通过`--v` 指定当前包的版本。
4.  你也可以在`$Home/.zctl/config.yaml `下指定这些配置。

### 创建一个GRPC服务greeter-service
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
### 创建一个调用服务greeter-api
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

## Golang项目
### 创建一个项目
创建一个目录并创建go.mod文件
```
module github.com/SunMaybo/metadata
```
![img_2.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img_2.png)
### proto文件管理
![img.png](https://raw.githubusercontent.com/SunMaybo/zero/develop/img/img9.png)
### 生成greeter-service服务
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
## 数据校验
## 授权机制
## 熔断降级
## 链路追踪



