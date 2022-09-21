package framework

// 服务提供者 (模块) 必要的对外接口
// 提供注册入容器的服务的接口, 又外部需要把自己注册入容器的服务去实现这些接口
// 这样就可以调用容器接口把自己注册进容器了
//- 获取服务凭证的能力 Name；
//- 创建服务实例化方法(模块实例化)的能力 Register；
//- 获取服务实例化方法参数的能力 Params；
//- 两个与实例化控制相关的方法，
//- 控制实例化时机方法 IsDefer
//- 实例化预处理的方法 Boot。

type NewInstance func(...any) (any, error)

// ServiceProvider 定义一个服务提供者需要实现的接口
type ServiceProvider interface {
	// Register 在服务容器中注册了一个实例化服务的方法，是否在注册的时候就实例化这个服务，需要参考 IsDefer 接口。
	Register(Container) NewInstance
	// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。
	// 如果 Boot 返回 error，整个服务实例化就会实例化失败，返回错误
	Boot(Container) error
	// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次 make 的时候进行实例化操作
	// false 表示不需要延迟实例化，在注册的时候就实例化。true 表示延迟实例化
	IsDefer() bool
	// Params params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
	Params(Container) []any
	// Name 代表了这个服务提供者的凭证
	Name() string
}
