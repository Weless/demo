package infra

import "github.com/tietang/props/kvs"

// 应用程序的启动管理器
type BootApplication struct {
	conf kvs.ConfigSource
	starterContext StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication{
	b := &BootApplication{conf:conf,starterContext: StarterContext{}}
	b.starterContext[KeyPros] = conf
	return b
}

func (b *BootApplication)Start()  {
	// 1. 初始化starter
	b.init()
	// 2. 安装starter
	b.setup()
	// 3. 启动starter
	b.start()
}

func (b *BootApplication)init()  {
	for _,starter := range StarterRegister.AllStarters(){
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication)setup()  {
	for _,starter := range StarterRegister.AllStarters(){
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication)start()  {
	for i,starter := range StarterRegister.AllStarters(){
		if starter.StartBlocking(){
			// 最后一个starter正常启动
			if i+1 == len(StarterRegister.AllStarters()){
				starter.Start(b.starterContext)
			}else{
				go starter.Start(b.starterContext)
			}
		}else{
			starter.Start(b.starterContext)
		}
	}
}
