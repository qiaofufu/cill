package cill

import (
	"fmt"
	"path"
	"strings"
)

/*
	使用Trie树是为了实现动态路由
	支持两种模式：
		(1) :name
		(2) *
*/
type node struct {
	part        string
	isFinalNode bool
	isOptional   bool // if node is `:name` or `*` isOptional is true
 	children  []*node
	pattern string
}

func (n *node) isMatch(part string) bool {
	if n.part == part || n.isOptional {
		return true
	} 
	return false
}

func (n *node) getMatchChild(part string) *node {
	for _, child := range n.children {
		if child.isMatch(part) {
			return child
		}
	}
	return nil
}

// getMatchChildren 获取匹配part的所有子节点
func (n *node) getMatchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.isMatch(part) {
			children = append(children, child)
		}
	}
	return children
}

func (n *node) insert(pattern string, parts []string, k int) error {
	// check if n is a final node 
	if k == len(parts) {
		n.isFinalNode = true
		n.pattern = pattern
		return nil
	}

	// get match child node
	part := parts[k]
	if len(part) == 0 {
		return fmt.Errorf("error syntax pattern: %s, %d`st len is 0", path.Join(parts...), k)
	}
	child := n.getMatchChild(part)

	// `part` node is not exist, create this node 
	if child == nil {
		child = &node{part: part, isOptional: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	return child.insert(pattern, parts, k + 1)
}

// Search 获取匹配parts的final节点
func (n *node) Search(parts []string, k int) *node {
	// 当url parts到达末尾时，或当前node出现*时，进行判断是否为final节点  
	if k == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.isFinalNode {
			return n
		}
		return nil
	}

	part := parts[k]
	children := n.getMatchChildren(part)
	for _, v := range children {
		if res := v.Search(parts, k + 1); res != nil {
			return res
		} 
	}
	return nil
}
