package framework

import (
	"errors"
	"strings"
)

// 代表树结构
type Tree struct {
	root *node // 根节点
}

// 代表节点
type node struct {
	isLast   bool                // 该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	segment  string              // uri中的字符串
	handlers []ControllerHandler // 控制器
	childs   []*node             // 子节点
}

func newNode() *node {
	return &node{
		isLast:   false,
		segment:  "",
		childs:   []*node{},
		handlers: make([]ControllerHandler, 0),
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}

// 判断一个segment是否是通用segment，即以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	// n.childs为空直接返回
	// n wei nil
	if len(n.childs) == 0 {
		return nil
	}
	// 如果segment是通配符, 则所有子节点都满足需求
	if isWildSegment(segment) {
		return n.childs
	}
	// slice初始化的时候, 长度要为0, 否者第一个节点会是nil!!!
	nodes := make([]*node, 0, len(n.childs))
	// 过滤下一层所有的子节点
	for _, cnode := range n.childs {
		// 子节点中有有通配符, 满足需求
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			// 子节点没通配符, 但是文本完全匹配 满足
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点数中存在了
func (n *node) matchNode(uri string) *node {
	// 使用分隔符将uri切割为两个部分
	segments := strings.SplitN(uri, "/", 2)
	// 第一个部分用于匹配下一层子节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)
	// 如果当前子节点没有一个符合，那么说明这个uri一定是之前不存在, 直接返回nil
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	// 如果只有一个segment，则是最后一个标记
	if len(segments) == 1 {
		// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}

		// 都不是最后一个节点
		return nil
	}

	// 如果有2个segment, 递归每个子节点继续进行查找
	for _, tn := range cnodes {
		//if tn == nil {
		//	panic("Tn is nil")
		//}
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// AddRouter 增加路由节点

// /book/list
// /book/:id (冲突)
// /book/:id/name
// /book/:student/age
// /:user/name
// /:user/name/:age (冲突)

func (tree *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("router exist: " + uri)
	}
	// 分割uri
	segments := strings.Split(uri, "/")
	// 遍历uri分割后的节点 for
	for index, segment := range segments {
		// 非通配符一律转大写
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		// 改节点是否已经在树中存在
		var objNode *node // 标记是否有该segment的节点
		childNode := n.filterChildNodes(segment)
		// 存在-遍历子节点,
		if len(childNode) > 0 {
			for _, cnode := range childNode {
				// 如果有segment相同的子节点，则选择这个子节点
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}
		isLast := index == len(segments)-1
		// 不存在
		if objNode == nil {
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = append(cnode.handlers, handlers...)
			}
			// 将新节点加入到childs中
			n.childs = append(n.childs, cnode)
			objNode = cnode //记录新节点, 外部做偏移用
		}
		n = objNode
	}
	return nil
	// 没有就创建该节点, 然后继续下一轮
}

// 匹配uri
func (tree *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}
