package vin

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
		//如果找到的字节点的part与需要寻找的part相等，
		// 或者该节点开启了精确匹配，那么就返回该child
		if child.part == part || child.isWild {
			return child
		}
	}

	//如果没有找到符合条件的，就返回nil
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
	//分割的parts有几层，就在第几层赋值
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//如果层数与height不想等，那么就继续向下插入
	part := parts[height]
	child := n.matchChild(part)
	//如果没有在当前节点的子节点中，找到对应的节点，就新建一个节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//有了对应的子节点以后，就继续向下插入
	child.insert(pattern, parts, height + 1)

}

func (n *node) search(parts []string, height int) *node {
	//如果找到了最后一层，或者当前节点的part的第一个字符为 *，
	//因为找到了 : 的话，还要继续查找，而*表示从当前节点开始可以结束搜索了
	//且该节点的pattern不为空时，返回这个节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	//如果不满足条件，则开始向下寻找
	part := parts[height]
	//向该节点的子节点中，寻找满足条件的子节点
	children := n.matchChildren(part)
	for _, child := range children {
		//在满足节点的子节点中，继续寻找对应层满足条件的子节点
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	}

	//没有满足条件的则返回nil
	return nil
}


