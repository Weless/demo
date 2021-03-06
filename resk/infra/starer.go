package infra

import "github.com/tietang/props/kvs"

const (
	KeyPros = "_conf"
)

// 资源启动器上下文
// 用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type StarterContext map[string]interface{}

func (s StarterContext)Props() kvs.ConfigSource {
	p := s[KeyPros]
	if p == nil{
		panic("配置文件还没初始化")
	}
	return p.(kvs.ConfigSource)
}

//资源启动器，每个应用少不了依赖其他资源，比如数据库，缓存，消息中间件等等服务
//启动器实现类，不需要实现所有方法，只需要实现对应的阶段方法即可，可以嵌入@BaseStarter
//通过实现资源启动器接口和资源启动注册器，友好的管理这些资源的初始化、安装、启动和停止。
//Starter对象注册器，所有需要在系统启动时需要实例化和运行的逻辑，都可以实现此接口
//注意只有Start方法才能被阻塞，如果是阻塞Start()，同时StartBlocking()要返回true
type Starter interface {
	//资源初始化和，通常把一些准备资源放在这里运行
	Init(StarterContext)
	//资源的安装，所有启动需要的具备条件，使得资源达到可以启动的就备状态
	Setup(StarterContext)
	//启动资源，达到可以使用的状态
	Start(StarterContext)
	//说明该资源启动器开始启动服务时，是否会阻塞
	//如果存在多个阻塞启动器时，只有最后一个阻塞，之前的会通过goroutine来异步启动
	//所以，需要规划好启动器注册顺序
	StartBlocking() bool
	//资源停止：
	// 通常在启动时遇到异常时或者启用远程管理时，用于释放资源和终止资源的使用，
	// 通常要优雅的释放，等待正在进行的任务继续，但不再接受新的任务
	Stop(StarterContext)
}

//服务启动注册器
//不用需外部构造，全局只有一个
type starterRegister struct {
	starter []Starter
}

var StarterRegister *starterRegister = &starterRegister{}

//返回所有的启动器
func (r *starterRegister) AllStarters() []Starter {
	return r.starter
}

//注册启动器
func (r *starterRegister) Register(starter Starter) {
	r.starter = append(r.starter,starter)
}

//注册starter
func Register(starter Starter) {
	StarterRegister.Register(starter)
}

//获取所有注册的starter
func GetStarters() []Starter {
	return StarterRegister.AllStarters()
}


//默认的空实现,方便资源启动器的实现
type BaseStarter struct {
}

func (s *BaseStarter) Init(ctx StarterContext)      {}
func (s *BaseStarter) Setup(ctx StarterContext)     {}
func (s *BaseStarter) Start(ctx StarterContext)     {}
func (s *BaseStarter) Stop(ctx StarterContext)      {}
func (s *BaseStarter) StartBlocking() bool          { return false }

