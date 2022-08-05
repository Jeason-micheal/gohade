package framework

// IGroup 定义IGroup, 做通用前缀匹配
// 这样在变更需求的时候, 只要修改或者重复实现该接口就好了
type IGroup interface {
	// 实现HttpMethod方法
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// 批量设置中间件
	Use(middlewares ...ControllerHandler)
	// 实现嵌套
	// 接口返回其自己
	Group(uri string) IGroup
}

// Group 实现接口的结构
type Group struct {
	core        *Core  // 指向core结构
	parent      *Group // 记录父组, 实现嵌套
	prefix      string // 这个group的通用前缀
	middlewares []ControllerHandler
}

// 初始化Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

// 实现Get方法
func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.middlewares, handlers...)
	g.core.Get(uri, allHandlers...)
}

// 实现Post方法
func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getAllHandlers(), handlers...)
	g.core.Post(uri, allHandlers...)
}

// 实现Put方法
func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getAllHandlers(), handlers...)
	g.core.Put(uri, allHandlers...)
}

// 实现Delete方法
func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getAllHandlers(), handlers...)
	g.core.Delete(uri, allHandlers...)
}

// 获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.getAbsolutePrefix() + g.prefix
}

func (g *Group) getAllHandlers() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.getAllHandlers(), g.middlewares...)
}

// 实现 Group 方法
func (g *Group) Group(uri string) IGroup {
	pgroup := NewGroup(g.core, uri)
	pgroup.parent = g
	return pgroup
}

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}
