# Min-Go

## 项目初始化

min-go 支持`go mod`初始化

```
require (
	github.com/duanchi/min/v2 v1.6.14
)
```



## 配置文件

### 配置结构体

项目中所有的配置可以通过定义一个`struct`类型的变量, 并由config/application.yaml填充配置项内容, 该变量可以引用`github.com/duanchi/min/v2/types`下的`Config`结构体进行组合。

```go
package config

import "github.com/duanchi/min/v2/types"

var Config = struct {
	types.Config	`yaml:",inline"`

	Jwt struct {
		SignatureKey string `yaml:"jwt_signature_key" default:"${JWT_SIGNATURE_KEY}"`
		ExpiresIn  int64 `yaml:"jwt_expires_in" default:"${JWT_EXPIRES_IN:7200}"`
	} // 自定义配置

	Beans struct {} // IOC容器初始化
}{}
```

在配置结构体中, 需要定义`yaml`的标签指明扩展配置文件的对应位置, 也可以通过`default`标签来指定配置的默认值

> 配置文件、默认值中可以使用`${}`的方式指定对应的环境变量, 也可以使用`${ENV_KEY:default_value}`的方式指定当未找到该环境变量时赋予的默认值

> 配置字段中的实际赋值可根据字段基本类型自动进行类型转换, 支持的基本类型有`int/int64` `float64` `string` `bool`

> 配置yaml文件默认配置在`config/application.yaml`, 可以通过Bootstrap方法更改配置文件位置

```yaml
env: "${ENV:development}"
#服务配置
application:
  server_port: "${SERVER_PORT:9801}"
  jwt_signature_key: "${JWT_SIGNATURE_KEY}"
  jwt_expire_in: "${JWT_EXPIRE_IN:7200}"
db:
  enabled: true
  dsn: "${DB_DSN:postgres://postgres:shvEVodOcqqTCWJ0@61.55.158.34:58932/cloud?sslmode=disable&prefix=dp_}"
```



```go
func main() {
	heron.SetConfigFile("./config/dashboard.yaml")
	heron.Bootstrap(&Config)
	return
}
```



### 获取配置值

在项目初始化后，可以直接引入该变量读取初始化后的配置。也可以通过配置获取方法`config.Get`获取配置

#### 直接获取配置

```go
import xxx/config
Dsn := Config.Db.Dsn
// host=172.31.16.1 port=3308 user=tb_cloud password=123456 dbname=thingsboard sslmode=disable
```

> 推荐使用配置变量直接获取配置, 既可以准确定位配置元素, 又可以省去类型强制转换。



#### 通过Get方法获取配置

```go
package config

func Get(key string) interface{} {}
```

其中`key`是以`.`分割的配置定义层级，获取配置后，需要进行强制类型转换。

```go
import "heurd.com/config"

func getConfig () {
    fmt.Println(config.Get("Db.Enabled").(bool))
    fmt.Println(config.Get("Db.Dsn").(string))
}

// true
// host=172.31.16.1 port=3308 user=tb_cloud password=123456 dbname=thingsboard sslmode=disable
```

> 读取配置后, 需手动进行类型转换



#### 在Bean中注入配置值

在Bean(可参考IOC容器章节)中, 可以使用`value`标签以类似配置文件中获取环境变量的方式获取配置值

```go
struct Sample struct {
  Expires string `value:"${Jwt.Expires}"`
}
```

> 为保证值可以正确注入, 需要将被注入的属性或字段设置为`可导出`的

## 项目文件结构

建议采用三层分离的文件结构，即`控制器`、`业务逻辑服务`、`实体关系映射`三层，分别对应`controller`、`sevice`、`mapper`三个包。

root

- main.go
- config.go
- controller
  - xxx.go `控制器文件`
- service
  - xxx.go `业务逻辑文件`
- mapper
  - xxx.go `数据库mapper文件`

## IOC容器

借鉴于`Java Spring`的`IOC容器`概念，可以通过定义`Bean`进行项目中实体实例的管理，在初始化后可在任何位置调用。

### Bean定义

使用类似配置文件定义时使用的自定义结构体进行`Bean`的定义，字段名为`Bean`的`name`, 字段类型为`Bean`实例的类型， 可通过标签配置扩展`Bean`配置信息

```go
package config

var Beans = struct {
	DataDevices controller.DevicesController `route:"data/devices"`
  [bean name] [bean struct] `[bean tags]`
}{}
```

若需实现上述Bean的IOC/DI，需继承`min.abstract.Bean`类型，并实现 `Init`方法

> `min.abstract.Bean`已经实现了空的`Init`方法，若初始化时不需任何额外操作，可不进行方法重写。



> 预定义的抽象基础类`min.abstract.RestController`和 `min.abstract.Service`已经继承了`min.abstract.Bean`类型，可以直接使用并在`Bean`定义文件中引用

```go
type Devices struct {
	abstract.Service
	Config string `value:"${Db.Dsn,172.31.128.5}"`
	Data *data.Devices `autowired:"true"`
}

func (service *Devices) Init() {
	fmt.Println("Inited!") // 初始化时将打印 'Inited!'
}
```



### Bean初始化

Bean定义完成后, 通过加入到`Config`文件的 `Beans`结构中, 可以实现Bean的自动初始化

```go
var Config = struct {
	types.Config	`yaml:",inline"`
  // ...
	Beans struct {

		GatewayMiddleware middleware.GatewayMiddleware `middleware:"true"`

		AuthorizationToken authorization.TokenController `rest:"authorization/token"`

		AccountService authorization2.AccountService
		TokenService authorization2.TokenService
	}
}{}
```



#### Bean初始化切面

可通过实现`Bean`类型的`Init`方法实现`Bean`初始化后的操作

#### [bean name]

用于标识当前`Bean`的`name`，可在min.GetBean方法参数中使用

#### [bean struct]

用于当前`Bean`类型或类的初始化，初始化后，将生成一个当前类型的实例

#### [bean tags]

用于描述当前`Bean`的特定属性，可在初始化后进行扩展操作

`bean tag`的可取值如下

##### `value`

用于设置当前属性或字段的值，可以使用`${[config-stack]}`来进行配置文件的读取，读取的值将自动转换为当前字段的类型，支持的类型有`int`、`int64`、`float64`、`string`、`bool`、`struct`

##### `autowired`

用于在`Bean`初始化后，将含有当前字段类型的`Bean`自动装载至当前的字段中，装载时将装载对象的引用类型。

> 自动装载时，字段类型需设置为将要装载类型`Bean`的指针类型，并将`autowired`值设置为`true`

##### `route`、`rest`、`method`

详见`Http服务`章节中的`路由配置`

##### `middleware`

用于定义基于HTTP服务的中间件处理方法, 如定义一个登录或者令牌授权认证、请求日志记录等, 通过继承`abstract.middlewire`定义一个中间件处理方法

```go
type AuthorizationMiddleware struct {
	abstract.Middleware
	TokenService *service.TokenService `autowired:"true"`
}

func (this *AuthorizationMiddleware) AfterRoute (ctx *gin.Context) {
  // 具体处理方法
	ctx.Next()
	return
}
```



#####  `扩展取值`

提供可自定义的`tag`标签，并通过扩展方法执行`Bean`初始化时的扩展操作。

可通过继承`types.BeanParser`来实现一个自定义的Bean处理扩展, 类似Java Spring中的注解定义

```go
type NativeApiBeanParser struct {
	types.BeanParser
}

func (parser NativeApiBeanParser) Parse (tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	resource := tag.Get("native_api")

	if resource != "" {
		NativeApiBeans[resource] = bean
	}
}
```

在Config结构体中加入配置

```
var Config = struct {
	config.Config	`yaml:",inline"`
	// ...
	Beans struct {
		NativeApiRoutesController native_api.RoutesController `native_api:"routes"`
		NativeApiServicesController native_api.ServicesController `native_api:"services"`
	}
}{
	Config: config.Config{
		Config: types.Config{
			BeanParsers: []_interface.BeanParserInterface{
				&native_api.NativeApiBeanParser{},
			},
		},
	},
}
```



项目初始化后可通过`min.GetBean`方法获取`Bean`

```go
import min

func Test() {
  min.GetBean("Fetch").(Fetch).Get()
}

```



> **只有在`Config`中定义的Bean才可以使用Bean的各种特性, 如自动初始化、全局对象、自动注入等特性**

## Http服务

项目初始化后，可以使用`min.HttpServer`进行Http服务器的相关操作

服务器端口配置可以使用配置`ServerPort` 指定

> Http服务器使用[Gin](https://github.com/gin-gonic/gin)

### 路由配置

请求路由目前可根据资源进行路由配置，是通过扩展`Bean`定义实现的，因此只需在`Bean`变量中添加名为`route`或`rest`或`restful`的标签, 并设置其属性值为绑定路径或资源路径即可，`Bean`字段类型为相应的处理结构体。

另外还可以追加使用`method`方法限定该路由处理方法可以处理的HTTP请求类型, 多个类型以`,`分割

> `rest`路由配置完成后，可自动配置`[resource]`, `[resource]/:id`和`[resource]/`的路由解析

### Restful控制器

创建Restful控制器只需新建一个结构体，组合`abstract.RestController`，并实现当前控制器关心的处理方法即可。

> 若不组合`abstract.RestController`，则需实现所有`min.interface.RestControllerInterface`方法

可实现的方法有

#### Fetch(获取) - GET

```go
func (controller RestController)
Fetch (
    id string,
    resource string,
    parameters *gin.Params,
    ctx *gin.Context
) (result interface{}, err types.Error) {}
```



#### Create(创建) - POST

```go
func (controller RestController)
Create (
    id string,
    resource string,
    parameters *gin.Params,
    ctx *gin.Context
) (result interface{}, err types.Error) {}
```



#### Update(更新) - PUT

```go
func (controller RestController)
Update (
    id string,
    resource string,
    parameters *gin.Params,
    ctx *gin.Context
) (result interface{}, err types.Error) {}
```



#### Remove(删除) - DELETE

```go
func (controller RestController)
Remove (
    id string,
    resource string,
    parameters *gin.Params,
    ctx *gin.Context
) (result interface{}, err types.Error) {}
```



#### Communicate(通信-基于websocket) - WEBSOCKET

```go
func (controller RestController)
Communicate (
    connection *websocket.Conn,
  	id string, resource string,
  	parameters *gin.Params,
  	ctx *gin.Context
) (err types.Error) {}
```



> 方法返回值可以是任意可转换为`json`的数据格式

> 其余方法未完成

## 数据库

若设置配置 `Db.Enabled=true`，则可开始数据库支持。在项目中可以使用`min.Db`进行和数据库有关的操作。

> 数据库使用[xorm](https://github.com/go-xorm/xorm)

## 异常

人为抛出异常时，建议使用`min.types.Error`类型，并通过`Message`和`Code`字段定义异常

Code需使用HTTP协议支持的错误码。

## 运行

main.go 必须包含以下结构

```go
min.Bootstrap(&Config, &Beans)
```

其中, `Config`为实际定义配置的变量，`Beans`为实际定义`Bean`的变量。项目初始化完成后，可通过引用调用配置和`Bean`的实际值。

项目开发环境可以使用 [fresh](https://github.com/gravityblast/fresh) 等hot-reload热加载方案。