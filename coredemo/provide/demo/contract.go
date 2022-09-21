package demo

// 存放服务的接口文件和服务凭证

const Key = "hade:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}