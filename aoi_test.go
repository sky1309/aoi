package aoi

import (
	"fmt"
	"testing"
)

type TestObject struct {
	node *Node
}

func (o *TestObject) OnEnter(node *Node) {
	fmt.Printf("my:(%f, %f), enter:(%f, %f)\n", o.node.x, o.node.y, node.x, node.y)
}

func (o *TestObject) OnMove(node *Node) {
	fmt.Printf("my:(%f, %f), move:(%f, %f)\n", o.node.x, o.node.y, node.x, node.y)
}

func (o *TestObject) OnLeave(node *Node) {
	fmt.Printf("my:(%f, %f), leave:(%f, %f)\n", o.node.x, o.node.y, node.x, node.y)
}

func newObjectNode(x, y, dis float32) *Node {
	node := NewNode(x, y, dis)

	testObj := &TestObject{node: node}
	node.SetAOIEvent(testObj)

	return node
}

func TestAOI(t *testing.T) {
	aoi := NewAOIManager()

	var nodes []*Node
	nodes = append(nodes, newObjectNode(1, 2, 3))
	nodes = append(nodes, newObjectNode(2, 2, 3))
	nodes = append(nodes, newObjectNode(8, 1, 3))
	nodes = append(nodes, newObjectNode(8, 1, 3))
	nodes = append(nodes, newObjectNode(7, 4, 3))

	for _, node := range nodes {
		aoi.Enter(node)
		fmt.Println("=====")
	}

	fmt.Println("------")
	for _, node := range nodes {
		aoi.Leave(node)
		fmt.Println("=====")
	}
}
