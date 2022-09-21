package framework

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// Bind 绑定一个服务提供者, 如果关键字凭证存在则替换, 返回error
	Bind(provide ServiceProvider) error
	// IsBind 检测是否已经绑定
	IsBind(key string) bool
	// Make 根据关键字凭证获取一个服务
	Make(key string) (any, error)
	// MustMake 根据关键字凭证凭证获取一个服务
	// 如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候
	// 请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) any
	// MakeNew 根据关键字凭证获取一个服务,且改服务不是单例模式的
	// 是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params ...any) (any, error)
}

type HadeContainer struct {
	Container // 组合 -- 结构体中直接放接口, 编译器就会检查是否有实现该接口
	// providers 存储注册的服务提供者, key是凭证
	providers map[string]ServiceProvider
	// 存放具体的实例, key是凭证
	instances map[string]any
	// Container中,注册事件比获取事件少很多, 属于读多写少, 读写锁更合适
	lock sync.RWMutex
}

func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]any{},
		lock:      sync.RWMutex{},
	}
}

func (hade *HadeContainer) PrintProvide() []string {
	ret := []string{}
	for _, provider := range hade.providers {
		name := provider.Name()
		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

func (hade *HadeContainer) Bind(provide ServiceProvider) error {
	hade.lock.Lock()
	defer hade.lock.Unlock()

	key := provide.Name()
	hade.providers[key] = provide
	// 不是延迟初始化
	if provide.IsDefer() == false {
		if err := provide.Boot(hade); err != nil {
			return err
		}
		// 实例化
		params := provide.Params(hade)
		method := provide.Register(hade)
		ins, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		hade.instances[key] = ins
	}
	return nil
}

func (hade *HadeContainer) IsBind(key string) bool {
	return hade.findServiceProvider(key) != nil
}

// Make 根据关键字凭证获取一个服务
func (hade *HadeContainer) Make(key string) (any, error) {
	return hade.make(key, nil, false)
}

// MustMake 根据关键字凭证凭证获取一个服务
// 如果这个关键字凭证未绑定服务提供者，那么会 panic。
// 所以在使用这个接口的时候
// 请保证服务容器已经为这个关键字凭证绑定了服务提供者。
func (hade *HadeContainer) MustMake(key string) any {
	ins, err := hade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return ins
}

// MakeNew 根据关键字凭证获取一个服务,且改服务不是单例模式的
// 是根据服务提供者注册的启动函数和传递的params参数实例化出来的
// 这个函数在需要为不同参数启动不同实例的时候非常有用
func (hade *HadeContainer) MakeNew(key string, params ...any) (any, error) {
	return hade.make(key, params, true)
}

// 真正的实例化一个服务
func (hade *HadeContainer) make(key string, params []any, forceNew bool) (any, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	// 查询是否已经注册了该服务提供者, 没有返回错误
	sp := hade.findServiceProvider(key)
	if sp == nil {
		return nil, fmt.Errorf("contract " + key + " have not register")
	}
	// 如果是new的 直接返回newInstance
	if forceNew {
		return hade.newInstance(sp, params)
	}
	//不需要重新实例化, 如果容器中已经实例化了, 就直接使用容器中的实例
	if ins, ok := hade.instances[key]; ok {
		return ins, nil
	}
	// 容器中没有实例化 则进行一次实例化
	ins, err := hade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	hade.instances[key] = ins
	return ins, nil
}

func (hade *HadeContainer) findServiceProvider(key string) ServiceProvider {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	if sp, ok := hade.providers[key]; ok {
		return sp
	}
	return nil
}

func (hade *HadeContainer) newInstance(sp ServiceProvider, params []any) (any, error) {
	if err := sp.Boot(hade); err != nil {
		return nil, err
	}
	// 是否需要原来sp中的参数??
	// 策略问题
	if params == nil {
		params = sp.Params(hade)
	}
	method := sp.Register(hade)

	// 正常的
	//ins, err := method(params...)
	//if err != nil {
	//	return nil, errors.New(err.Error())
	//}
	//return ins, err

	return method(params...)
}
