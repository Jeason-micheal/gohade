package framework

import "testing"

func Test_filterChildNodes(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(*Context) error {
			return nil
		},
		childs: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: func(c *Context) error {
					return nil
				},
				childs: nil,
			},
			{
				isLast:  false,
				segment: ":id",
				handler: nil,
				childs:  nil,
			},
		},
	}

	{
		nodes := root.filterChildNodes("FOO")
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}
	{
		nodes := root.filterChildNodes(":foo")
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}
}

func Test_matchNode(t *testing.T) {
	// func (n *node) matchNode(uri string) *node
	// 创建节点
	// 输入url 输出对应的节点
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(c *Context) error { return nil },
		childs: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: nil,
				childs: []*node{
					&node{
						isLast:  true,
						segment: "BAR",
						handler: func(c *Context) error { panic("not implemented") },
						childs:  []*node{},
					},
				},
			},
			{
				isLast:  false,
				segment: ":id",
				handler: nil,
				childs:  nil,
			},
		},
	}

	{
		node := root.matchNode("foo/bar")
		if node == nil {
			t.Error("match normal node error")
		}
	}

	{
		node := root.matchNode("test")
		if node == nil {
			t.Error("match test")
		}
	}

}