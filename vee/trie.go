package vee

import "strings"

//该文件用来以前缀树来存储动态路由的信息
type node struct {
	pattern string //待匹配的路由
	part string  //路由中的一部分，一个api部分
	children []*node //该节点的子节点
	isWild bool //是否精确匹配，如果part中含有: 或者 *，就会对其精确匹配，此时为true,即动态路由的实现
}

//该函数是找到能够成功匹配的第一个节点，用于插入
func (n *node) matchChild(part string) *node {
	//在子节点中寻找与part相同的子节点，或者开启了精确匹配的子节点
	for _, child := range n.children {
		//找到一个符合条件的，就直接返回
		if child.part == part || child.isWild {
			return child
		}
	}
	//如果在子节点中没有找到匹配的，就返回nil
	return nil
}

//找到所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//用于在前缀树中插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	//只有在与parts的长度对应的层数时，才会给当前的节点赋值pattern
	//第0层为对应的方法，诸如GET，POST等，然后是第一层的API
	//例如/hello，会在第一层赋值，/hello/v1,会在第二层的节点赋值
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//在当前的节点的子节点里，寻找是否有满足条件的节点
	part := parts[height]
	child := n.matchChild(part)
	//如果没有找到，那么就创建一个节点
	if child == nil {
		//创建新的节点，如果part的首个字符为：或者*，那么就将isWild置为true
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		//并将创建的节点添加到当前节点的子节点中
		n.children = append(n.children, child)
	}
	//因为还未匹配结束，还需要继续插入
	child.insert(pattern, parts, height + 1)
}

func (n *node) search(parts []string, height int) *node {
	//如果找到了最后一层，或者当前节点的part的第一个字符为 *，
	//且该节点的patter不为空时，返回这个节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	//如果不满足条件，则继续向下寻找
	part := parts[height]
	//首先找到满足条件的子节点
	children := n.matchChildren(part)
	for _, child := range children {
		//如果子节点中有满足条件的节点，那么就返回满足的节点
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	}

	//没有满足条件的则返回nil
	return nil
}


