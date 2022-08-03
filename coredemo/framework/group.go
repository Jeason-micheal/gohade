package framework

// IGroup 定义IGroup, 做通用前缀匹配
// 这样在变更需求的时候, 只要修改或者重复实现该接口就好了
type IGroup interface {
	// 实现HttpMethod方法
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	// 实现嵌套
	// 接口返回其自己
	Group(uri string) IGroup
}

// Group 实现接口的结构
type Group struct {
	core   *Core  // 指向core结构
	parent *Group // 记录父组, 实现嵌套
	prefix string // 这个group的通用前缀
}

// 初始化Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

// 实现Get方法
func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Get(uri, handler)
}

// 实现Post方法
func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Post(uri, handler)
}

// 实现Put方法
func (g *Group) Put(uri string, handler ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Put(uri, handler)
}

// 实现Delete方法
func (g *Group) Delete(uri string, handler ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Delete(uri, handler)
}

// 获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.getAbsolutePrefix() + g.prefix
}

// 实现 Group 方法
func (g *Group) Group(uri string) IGroup {
	pgroup := NewGroup(g.core, uri)
	pgroup.parent = g
	return pgroup
}
