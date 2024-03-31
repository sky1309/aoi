package aoi

import (
	"math"
)

type AOIEvent interface {
	OnEnter(node *Node)
	OnMove(node *Node)
	OnLeave(node *Node)
}

func NewNode(x, y, dis float32) *Node {
	return &Node{
		x:   x,
		y:   y,
		dis: dis,
	}
}

type Node struct {
	x   float32
	y   float32
	dis float32 // sight range

	xPrev *Node
	xNext *Node
	yPrev *Node
	yNext *Node

	aoiEvent AOIEvent
}

func (n *Node) SetDis(dis float32) {
	n.dis = dis
}

func (n *Node) SetAOIEvent(aoiEvent AOIEvent) {
	n.aoiEvent = aoiEvent
}

// ----------------
// ----------------

func NewAOIManager() *AOIManager {
	aoiManager := &AOIManager{
		xHead: &Node{x: -1, y: -1},
		yHead: &Node{x: -1, y: -1},
	}
	return aoiManager
}

type AOIManager struct {
	xHead *Node
	yHead *Node
}

func (aoi *AOIManager) add(node *Node) {

	// add x
	find := false
	var prev *Node
	cur := aoi.xHead
	for cur != nil {
		if node.x < cur.x {
			find = true

			node.xPrev = cur.xPrev
			node.xNext = cur
			cur.xPrev.xNext = node
			cur.xPrev = node
			break
		}
		prev = cur
		cur = cur.xNext
	}
	if !find {
		prev.xNext = node
		node.xPrev = prev
	}

	// add y
	prev = nil
	find = false
	cur = aoi.yHead
	for cur != nil {
		if node.y < cur.y {
			find = true

			node.yPrev = cur.yPrev
			node.yNext = cur
			cur.yPrev.yNext = node
			cur.yPrev = node
			break
		}
		prev = cur
		cur = cur.xNext
	}

	if !find {
		prev.yNext = node
		node.yPrev = prev
	}
}

func (aoi *AOIManager) remove(node *Node) {
	// remove x
	node.xPrev.xNext = node.xNext
	if node.xNext != nil {
		node.xNext.xPrev = node.xPrev
	}
	node.xPrev = nil
	node.xNext = nil

	// remove y
	node.yPrev.yNext = node.yNext
	if node.yNext != nil {
		node.yNext.yPrev = node.yPrev
	}
	node.yPrev = nil
	node.yNext = nil
}

func (aoi *AOIManager) findNears(node *Node, dis float32) map[*Node]struct{} {
	nears := make(map[*Node]struct{})

	for cur := node.xPrev; cur != nil; cur = cur.xPrev {
		if cur == aoi.xHead {
			break
		}
		if node.x-cur.x > dis {
			break
		}
		if math.Abs(float64(cur.y-node.y)) <= float64(dis) {
			nears[cur] = struct{}{}
		}

	}

	for cur := node.xNext; cur != nil; cur = cur.xNext {
		// over distance
		if cur.x-node.x > dis {
			break
		}
		if math.Abs(float64(cur.y-node.y)) <= float64(dis) {
			nears[cur] = struct{}{}
		}
	}

	return nears
}

func (aoi *AOIManager) FindNears(node *Node, dis float32) map[*Node]struct{} {
	return aoi.findNears(node, dis)
}

func (aoi *AOIManager) Enter(node *Node) {
	aoi.add(node)

	nears := aoi.findNears(node, node.dis)
	for n := range nears {
		n.aoiEvent.OnEnter(node)
		node.aoiEvent.OnEnter(n)
	}
}

func (aoi *AOIManager) Leave(node *Node) {
	nears := aoi.findNears(node, node.dis)

	for n := range nears {
		n.aoiEvent.OnLeave(node)
		node.aoiEvent.OnLeave(n)
	}

	aoi.remove(node)
}

func (aoi *AOIManager) Move(node *Node, x, y float32) {
	beforeNears := aoi.findNears(node, node.dis)
	aoi.remove(node)

	node.x = x
	node.y = y
	aoi.add(node)
	newNears := aoi.findNears(node, node.dis)

	for n := range beforeNears {
		if _, ok := newNears[n]; ok {
			n.aoiEvent.OnMove(node)
		} else {
			n.aoiEvent.OnLeave(node)
			node.aoiEvent.OnLeave(n)
		}
	}

	for n := range newNears {
		if _, ok := beforeNears[n]; !ok {
			n.aoiEvent.OnEnter(node)
			node.aoiEvent.OnEnter(n)
		}
	}
}
