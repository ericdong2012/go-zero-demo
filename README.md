
## 对应视频讲解地址
> https://www.bilibili.com/video/BV1gu411q7S7?spm_id_from=333.999.0.0

## 0. why go-zero

> https://github.com/zeromicro/go-zero/blob/master/readme-cn.md
>
> https://go-zero.dev/cn/
>
> https://github.com/zeromicro/zero-examples/



### 1. 需求

```
首先，能够写更少代码解决业务需求。更少的代码意味着更快的产出，更少的bug。

其次，框架是否稳定，有没经过实战检验。

再次，社区是否活跃，遇到问题是否能够快速得到解决。

最后，发展怎么样。有没有上cncf, github star数怎么样，代码更新怎么样。
```



### 2. 特点

```
go-zero（收录于 CNCF 云原生技术全景图：https://landscape.cncf.io/?selected=go-zero）是一个集成了各种工程实践的 web 和 rpc 框架。通过弹性设计保障了大并发服务端的稳定性，经受了充分的实战检验。

go-zero 包含极简的 API 定义和生成工具 goctl，可以根据定义的 api 文件一键生成 Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript 代码，并可直接运行。

特点：
    强大的工具支持，尽可能少的代码编写
    极简的接口
    完全兼容 net/http
    面向故障编程，弹性设计
    内建级联超时控制、限流、自适应熔断、自适应降载等微服务治理能力，无需配置和额外代码
    内建服务发现、负载均衡
    内建自动缓存控制
    微服务治理中间件可无缝集成到其它现有框架使用
    大量微服务治理和并发工具包
    API 参数自动校验
    链路跟踪、监控报警等
    高并发支撑，稳定保障了疫情期间每天的流量洪峰
```



![](架构图.png)



### 3. 对比

|              | go-zero                  | kratos                 | dubbo-go    | sofarpc | tars          |
| ------------ | ------------------------ | ---------------------- | ----------- | ------- | ------------- |
| 适用语言     | go / http，rpc(protobuf) | go/grpc,  grpc-gateway | go          | java    | go(IDL, tars) |
| 使用难度     | 简单                     | 有一定难度             | 类Java      | ×       | 有难度        |
| cncf         | yes                      | yes                    | yes         | yes     | yes           |
| stars        | 14.1k                    | 16.3k                  | 3.8k        | 3.5k    | 2.9k          |
| 开源时间     | 2020.8.7                 | 2019.1.10              |             |         |               |
| 文档         | yes                      | yes                    | yes         | ×       | ×             |
| 微服务工具包 | yes                      | yes                    | no(单纯rpc) | ×       | ×             |











## 1. 快速入门单体服务

### 1.  安装

> 要求安装了 go sdk,   go > 1.13,   配置了 goroot,  gopath

```
# 开启go mod, 配置goproxy
	go env set -w GO111MODULE=on GOPROXY=https://goproxy.cn/,direct 
```



### 2. 创建greet工程

```
# 创建文件夹，进入文件夹
mkdir go-zero-demo
cd ~/go-zero-demo

# 初始化go mod 文件
go mod init go-zero-demo

# 安装go-zero
go get -u github.com/zeromicro/go-zero

# 安装goctl
    # go < 1.16
    go get -u github.com/zeromicro/go-zero/tools/goctl

    # go >= 1.16
    go install github.com/zeromicro/go-zero/tools/goctl@latest

# goctl 生成 greet 单体项目目录(golang会有乱码，可更改)
goctl api new greet
// 输出 ?[32mDone.?[0m
```



### 3. 编写业务逻辑

```
# greet/internal/logic/greetlogic.go

func (l *GreetLogic) Greet(req types.Request) (*types.Response, error) {
    return &types.Response{
        Message: "Hello go-zero",
    }, nil
}
```



### 4. 启动

```
go mod tidy
cd greet 
go run greet.go -f etc/greet-api.yaml
```



### 5. 测试

```
localhost:8888/from/you|me 
或者
curl -i -X GET http://localhost:8888/from/you
```





## 2. 快速开发微服务

### 1. 背景

```
本小节将以一个订单服务调用用户服务来简单演示一下，演示代码仅传递思路，其中有些环节不会一一列举。

订单服务(order) 提供一个查询接口
用户服务(user)  提供一个方法供订单服务获取用户信息


根据情景提要我们可以得知，订单是直接面向用户，通过http协议访问数据，而订单内部需要获取用户的一些基础数据，既然我们的服务是采用微服务的架构设计， 那么两个服务（user,order）就必须要进行数据交换，服务间的数据交换即服务间的通讯，到了这里，采用合理的通讯协议也是一个开发人员需要 考虑的事情，可以通过http，rpc等方式来进行通讯，这里我们选择rpc来实现服务间的通讯。从上文得知，我们需要一个

order api
user rpc
```



### 2. 安装protoc及插件

> 注意 gopath 在环境变量中

```
windows:

1. 在 https://github.com/protocolbuffers/protobuf/releases/tag/v3.3.0  下载protoc-3.3.0-win32.zip的包
2. 将protoc.exe 复制到gopath/bin下
3. 安装插件
# protoc-gen-go是用来将protobuf的的代码转换成go语言代码的一个插件
go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2

4. 测试
protoc --version

# 最好在cmd中运行，编辑器的terminal中不一定好用
protoc  --proto_path=IMPORT_PATH  --go_out=./  (--python_out=xxx)   *.proto

# 报错1： --go_out: protoc-gen-go: Plugin failed with status code 1.
protoc-gen-go 找不到
将 gopath/bin 加到环境变量就好了(实在不行 gopath/bin下执行就好了)
```



### 3. 创建mall工程及user rpc服务

> https://go-kratos.dev/docs/guide/api-protobuf

```
cd ~/go-zero-demo && go mod init go-zero-demo
mkdir -p mall/user/rpc && cd mall/user/rpc

# vim  user.proto

syntax = "proto3";

package user;

//protoc-gen-go 版本大于1.4.0, proto文件需要加上go_package,否则无法生成
option go_package = "./user";

message IdRequest {
  string id = 1;
}

message UserResponse {
// 用户id
string id = 1;
// 用户名称
string name = 2;
// 用户性别
string gender = 3;
}

service User {
  rpc getUser(IdRequest) returns(UserResponse);
}



cd mall/user/rpc && goctl rpc proto -src user.proto -dir .

注意事项：
新建userclient 文件夹，将原有user文件夹下的user.go 移动到 userclient中
（后续的其他rpc也可参考此方式，新建xxxclient, 然后移动文件）
```



### 4.  编写user业务逻辑

```
# internal/logic/getuserlogic.go


func (l *GetUserLogic) GetUser(in *user.IdRequest) (*user.UserResponse, error) {
	return &userclient.UserResponse{
		Id:   in.Id,
		Name: "test",
	}, nil
}
```



### 5. 创建order api服务

```
cd ~/go-zero-demo/mall
mkdir -p order/api && cd order/api

# 新建order.api
 type(
      OrderReq {
          Id string `path:"id"`
      }

      OrderReply {
          Id string `json:"id"`
          Name string `json:"name"`
      }
  )
  service order {
      @handler getOrder
      get /api/order/get/:id (OrderReq) returns (OrderReply)
  }


# 生成api 文件目录
cd mall/order/api && goctl api go -api order.api -dir .
```



### 6.  编写order业务逻辑

```
# 添加user rpc配置
# vim mall/order/api/internal/config/config.go

type Config struct {
	rest.RestConf
	UserRpc zrpc.RpcClientConf
}



# 添加yaml配置
# vim mall/order/api/etc/order.yaml

  Name: order
  Host: 0.0.0.0
  Port: 8888
  UserRpc:
    Etcd:
      Hosts:
      - 127.0.0.1:2379
      Key: user.rpc
      
      
      
# 完善服务依赖
# vim mall/order/api/internal/svc/servicecontext.go

  type ServiceContext struct {
      Config  config.Config
      UserRpc userclient.User
  }

  func NewServiceContext(c config.Config) *ServiceContext {
      return &ServiceContext{
          Config:  c,
          UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
      }
  }
  
  
# 添加业务逻辑
# vim mall/order/api/internal/logic/getorderlogic.go


func (l *GetOrderLogic) GetOrder(req types.OrderReq) (*types.OrderReply, error) {
  	user, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userclient.IdRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	if user.Name != "test" {
		return nil, errors.New("用户不存在")
	}

	return &types.OrderReply{
		Id:   req.Id,
		Name: "test order",
	}, nil
}
```



### 7. 测试

```
启动etcd
etcd

先启动user rpc
go mod tidy
go run user.go -f etc/user.yaml

后启动order api
go run order.go -f etc/order.yaml

访问:
http://localhost:8888/api/order/get/1
```



## 3. goctl

### 1. api

```
goctl api -h
   new       fast create api service
   format    format api files
   validate  validate api file
   doc       generate doc files
   go        generate go files for provided api in yaml file
   

new(快速创建一个api服务, 包括api文件和文件目录)
	示例：goctl api new user
	
-o(只生成api文件)
	示例：goctl api -o user.api

go(根据api文件生成代码)
    -dir (指定代码存放目录)
    -api (指定api文件源)
    -force (是否强制覆盖已经存在的文件)
    -style (指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)

	示例： goctl api go -api user.api -dir .
```

### 2. rpc

```
goctl rpc -h
   new       generate rpc demo service
   template  generate proto template
   proto     generate rpc from proto

# 生成标准的rpc目录代码
goctl rpc new greet

# 生成rpc proto
goctl rpc template -o=user.proto

# 根据rpc 文件生成代码
goctl rpc proto -src user.proto -dir .
```



### 3. model

```
1. 通过ddl生成
	 goctl model mysql ddl -src="./*.sql" -dir="./sql/model" -c
	 
	执行上述命令后即可快速生成CURD代码。
      model
      │   ├── error.go
      │   └── usermodel.go
      
2. 通过datasource生成
goctl model mysql datasource -url="root:mysql@tcp(127.0.0.1:3306)/gin-test" -table="*"  -dir="./model"
```



### 4. docker

```
-go(指定main函数文件)

-port(指定暴露端口)


goctl docker -go hello.go
```



### 5. kube

```
deploy
    -name 服务名称
    -namespace 指定k8s namespace
    -image 指定镜像名称
    -secret 指定获取镜像的k8s secret
    -requestCpu 指定cpu默认分配额
    -requestMem 指定内存默认分配额
    -limitCpu 指定cpu最大分配额
    -limitMem 指定内存最大分配额
    -o deployment.yaml输出目录
    -replicas 指定副本数
    -revisions 指定保留发布记录数
    -port 指定服务端口
    -nodePort 指定服务对外暴露端口
    -minReplicas 指定最小副本数
    -maxReplicas 指定最大副本数
    
goctl kube deploy -name greet -namespace greet-ns  -image greet:v1 -o greet.yaml -port 6379
```



### 6. template

> https://go-zero.dev/cn/template-cmd.html
>
> https://go-zero.dev/cn/template.html

```
1. goctl template init --home ./template
2. 更改对应的模板
3. goctl rpc new greet --home ./template
```





## 4. api语法

```
// api语法版本
syntax = "v1"

// import literal
import "foo.api"

// import group
import (
    "bar.api"
    "foo/bar.api"
)
info(
    author: "songmeizi"
    date:   "2020-01-08"
    desc:   "api语法示例及语法说明"
)

// type literal
type Foo{
    Foo int `json:"foo"`
}

// type group
type(
    Bar{
        Bar int `json:"bar"`
    }
)

// service block
@server(
    jwt:   Auth
    group: foo
)
service foo-api{
    @doc "foo"
    @handler foo
    post /foo (Foo) returns (Bar)
}
```

### 1. type

```
type:
    不支持alias
    不支持time.Time数据类型
    结构体名称、字段名称、不能为golang关键字   interface, var等都不行

type Foo{
    Id int `path:"id"`      // 大写开头，不能写成非导出字段，即小写开头
    Foo int `json:"foo"`
}

type Bar{
    Bar int `form:"bar"`
}

type(
    FooBar{
        FooBar int `json:"fooBar"`
    }
)
```

- tag值

  | key    | 描述                                                         | 提供方  | 有效范围          | 示例            |
  | ------ | ------------------------------------------------------------ | ------- | ----------------- | --------------- |
  | json   | json序列化tag                                                | golang  | request、response | `json:"fooo"`   |
  | path   | 路由path，如`/foo/:id`                                       | go-zero | request           | `path:"id"`     |
  | form   | 标志请求体是一个form（POST方法时）或者一个query(GET方法时`/search?name=keyword`) | go-zero | request           | `form:"name"`   |
  | header | HTTP header，如 `Name: value`                                | go-zero | request           | `header:"name"` |

- tag修饰符

  | key      | 描述                                  | 提供方  | 有效范围 | 示例                         |
  | -------- | ------------------------------------- | ------- | -------- | ---------------------------- |
  | optional | 定义当前字段为可选参数                | go-zero | request  | `json:"name,optional"`       |
  | options  | 定义当前字段的枚举值,多个以竖线\|隔开 | go-zero | request  | `json:"gender,options=male"` |
  | default  | 定义当前字段默认值                    | go-zero | request  | `json:"gender,default=male"` |
  | range    | 定义当前字段数值范围                  | go-zero | request  | `json:"age,range=[0:120]"`   |

  **tag修饰符需要在tag value后以引文 逗号, 隔开**



### 2. service

```


service foo-api{
  @doc "foo"
  @handler foo
  post /foo/:id (Foo) returns (Bar)
}

@server(
  jwt: Auth
  group: foo
  middleware: AuthMiddleware
  prefix: /api
)

service foo-api{
  @handler ping
  get /ping

  @doc "foo"
  @handler bar
  post /bar/:id (Foo)
}
```



@service 关键key描述说明

| key        | 描述                                                         | 示例                         |
| ---------- | ------------------------------------------------------------ | ---------------------------- |
| jwt        | 声明当前service下所有路由需要jwt鉴权，且会自动生成包含jwt逻辑的代码 | `jwt: Auth`                  |
| group      | 声明当前service或者路由文件分组                              | `group: login`               |
| middleware | 声明当前service需要开启中间件                                | `middleware: AuthMiddleware` |
| prefix     | 添加路由分组                                                 | `prefix: /api`               |



## 5. 其他

### 1. api配置

```
https://go-zero.dev/cn/api-config.html
```



### 2. rpc配置

```
https://go-zero.dev/cn/rpc-config.html
```



### 3. 业务开发

```
https://go-zero.dev/cn/service-design.html
```

#### 1. jwt

> https://talkgo.org/t/topic/1114

```
1. 新建jwt, jwt.api

type JwtTokenRequest struct {
}

type JwtTokenResponse struct {
  AccessToken  string `json:"access_token"`
  AccessExpire int64  `json:"access_expire"`
  RefreshAfter int64  `json:"refresh_after"` // 建议客户端刷新token的绝对时间, 告诉客户端什么时候该刷新	JWT token了，一般都需要设置过期时间前几天。
}

type GetUserRequest struct { 
  UserId string `json:"userId"`
}

type GetUserResponse struct {
  Name string `json:"name"`
}

service jwt-api {
  @server(
    handler: JwtHandler
  )
  post /user/token(JwtTokenRequest) returns (JwtTokenResponse)
}

@server(
  jwt: JwtAuth
)
service jwt-api {
  @server(
    handler: GetUserHandler
  )
  post /user/info(GetUserRequest) returns (GetUserResponse)
}


goctl api go -api jwt.api -dir .


2. 修改jwtlogic.go
func (l *JwtLogic) Jwt(req types.JwtTokenRequest) (*types.JwtTokenResponse, error) {
	var accessExpire = l.svcCtx.Config.JwtAuth.AccessExpire

	now := time.Now().Unix()
	accessToken, err := l.GenToken(now, l.svcCtx.Config.JwtAuth.AccessSecret, nil, accessExpire)
	if err != nil {
		return nil, err
	}

	return &types.JwtTokenResponse{
    AccessToken:  accessToken,
    AccessExpire: now + accessExpire,
    RefreshAfter: now + accessExpire/2,
  }, nil
}

func (l *JwtLogic) GenToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}




3. 更改etc/jwt-api.yaml

Name: jwt-api
Host: 0.0.0.0
Port: 8888
JwtAuth:
  AccessSecret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AccessExpire: 604800
  
  
4. 启动服务，访问
go run jwt.go -f etc/jwt-api.yaml

curl -X POST '127.0.0.1:8888/user/token'

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDMxMDU5NTgsImlhdCI6MTY0MjUwMTE1OH0.5GtLaOC58lCdFSkprNkG-mWEDgaf3yZ8xZkiLv9Dh5o",
    "access_expire": 1643105958,
    "refresh_after": 1642803558
}



5. 修改getuserlogic.go
func (l *GetUserLogic) GetUser(req types.GetUserRequest) (*types.GetUserResponse, error) {
	return &types.GetUserResponse{Name: "kim"}, nil
}


6. 重新测试
我们先不带JWT Authorization header请求头测试下，返回http status code是401，符合预期。
curl -X POST '127.0.0.1:8888/user/info' \
-H 'Content-Type: application/json' \
-d '{
    "userId": "a"
}'

//结果
status: 401 Unauthorized



加上Authorization header请求头测试。
curl -X POST '127.0.0.1:8888/user/info' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEyNjE0MjksImlhdCI6MTYwMDY1NjYyOX0.6u_hpE_4m5gcI90taJLZtvfekwUmjrbNJ-5saaDGeQc' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userId": "a"
}'

// 结果
{"name":"kim"}
status: 200 ok
```







#### 2. 中间件

> https://go-zero.dev/cn/middleware.html



**路由中间件**

```
1. 更改jwt.api
 
 @server(
	jwt: JwtAuth
	middleware: ExampleMiddleware
)

goctl api go -api jwt.api -dir .


	变化：
        1. 新增了middleware
        2. routes下 
                rest.WithMiddlewares(
                    []rest.Middleware{serverCtx.ExampleMiddleware},
			
			


2. 完善资源依赖ServiceContext
type ServiceContext struct {
	Config        config.Config
	ExampleMiddleware rest.Middleware    // 要求和yaml文件中定义的一致
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ExampleMiddleware: middleware.NewExampleMiddleware().Handle,
	}
}

3. 编写中间件逻辑
middleware/examplemiddleware.go

func (m *ExampleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
      return func(w http.ResponseWriter, r *http.Request) {
          // TODO generate middleware implement function, delete after code implementation

          // Passthrough to next handler if need
          logx.Info("router middleware")
          next(w, r)
      }
  }
  
  
4. 启动服务，访问
curl -X POST '127.0.0.1:8888/user/info' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEyNjE0MjksImlhdCI6MTYwMDY1NjYyOX0.6u_hpE_4m5gcI90taJLZtvfekwUmjrbNJ-5saaDGeQc' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userId": "a"
}'

// 结果
{"@timestamp":"2022-01-18T18:58:01.785+08","level":"info","content":"router middleware"}
{"@timestamp":"2022-01-18T18:58:01.786+08","level":"info","content":"[HTTP] POST - 200 - /user/info - 127.0.0.1:10114 - PostmanRuntime/7.26.8 - 0.4ms","trace":"63f6ee648676d619eddeac9427464d2b","span":"31c2755d13987c76"}
```



**全局中间件**

```
通过rest.Server提供的Use方法即可

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    ctx := svc.NewServiceContext(c)
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()

    // 全局中间件
    server.Use(func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            logx.Info("global middleware")
            next(w, r)
        }
    })


    handler.RegisterHandlers(server, ctx)

    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}

运行，访问
curl -X POST '127.0.0.1:8888/user/token'

// 结果
{"@timestamp":"2022-01-18T19:01:21.882+08","level":"info","content":"global middleware"}
{"@timestamp":"2022-01-18T19:01:21.883+08","level":"info","content":"[HTTP] POST - 200 - /user/token - 127.0.0.1:10123 - PostmanRuntime/7.26.8 - 0.5ms","trace":"7e30b33977a50fed8321e5b1f47e60bc","span":"d590ca672987276a"}

```



中间件中调用其他 中间件

```
通过闭包的方式把其它服务传递给中间件，示例如下：

// 模拟的其它服务
type AnotherService struct{}

func (s *AnotherService) GetToken() string {
    return stringx.Rand()
}

// 常规中间件
func middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("X-Middleware", "static-middleware")
        next(w, r)
    }
}

// 调用其它服务的中间件
func middlewareWithAnotherService(s *AnotherService) rest.Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Add("X-Middleware", s.GetToken())
            next(w, r)
        }
    }
}
```





### 4. ci/cd

```
https://go-zero.dev/cn/ci-cd.html

https://go-zero.dev/cn/service-deployment.html
```



### 5. 日志收集

```
https://go-zero.dev/cn/log-collection.html
```



### 6. 链路追踪

```
https://go-zero.dev/cn/trace.html
```



### 7. 监控管理

```
https://go-zero.dev/cn/service-monitor.html
```



### 8. 开发规范

#### 命名规范

**命名准则**

- 当变量名称在定义和最后一次使用之间的距离很短时，简短的名称看起来会更好。
- 变量命名应尽量描述其内容，而不是类型
- 常量命名应尽量描述其值，而不是如何使用这个值
- 在遇到for，if等循环或分支时，推荐单个字母命名来标识参数和返回值
- method、interface、type、package推荐使用单词命名
- package名称也是命名的一部分，请尽量将其利用起来
- 使用一致的命名风格

**文件命名规范**

- 全部小写
- 除unit test外 避免下划线(_)
- 文件名称不宜过长

**变量命名规范参考**

- 首字母小写
- 驼峰命名
- 见名知义，避免拼音替代英文
- 不建议包含下划线(_)
- 不建议包含数字

**函数、常量命名规范**

- 驼峰式命名
- 可exported的必须首字母大写
- 不可exported的必须首字母小写
- 避免全部大写与下划线(_)组合

#### 路由规范

- 推荐脊柱式命名
- 小写单词 + 横杠(-)组合
- 见名知义

  

```
/user/get-info
/user/get/info
/user/password/change/:id
```



#### 编码规范

**import**

- 单行import不建议用圆括号包裹

- 按照`官方包`，NEW LINE，`当前工程包`，NEW LINE，`第三方依赖包`顺序引入

  ```go
    import (
        "context"
        "string"
  
        "greet/user/internal/config"
  
        "google.golang.org/grpc"
    )
  ```

**函数返回**

- 对象避免非指针返回
- 遵循有正常值返回则一定无error，有error则一定无正常值返回的原则



**错误处理**

- 有error必须处理，如果不能处理就必须抛出。
- 避免下划线(_)接收error



**函数体编码**

- 建议一个block结束空一行，如if、for等

  ```go
    func main () {
        if x==1 {
            // do something
        }
  
        fmt.println("xxx")
    }
  ```

- return前空一行

  ```go
    func getUser(id string)(string,error){
        ....
        return "xx",nil
    }
  ```