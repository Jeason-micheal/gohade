package demo

import (
	"github.com/gohade/hade/framework"
)

// 实现具体的服务实例

type DemoService struct {
	// 接口实现
	Service
	// 参数
	c framework.Container
}

//  func(...any) (any, error)
// 其中 参数一定有 framework.Container
func NewDemoService(params ...any) (any, error) {
	// 展开参数
	var c framework.Container
	if len(params) > 0 {
		// 这里出发RUnlock 为什么
		c = params[0].(framework.Container)
	}
	//c := params[0].(framework.Container)
	//fmt.Println("new demo service")
	return &DemoService{c: c}, nil
}

func (s *DemoService) GetFoo() Foo {
	return Foo{Name: "i am foo"}
}
