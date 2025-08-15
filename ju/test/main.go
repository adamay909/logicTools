package main

import (
	"fmt"

	"github.com/adamay909/logicTools/ju"
)

type mynodeval interface {
	isNodeVal()
}

type mynode = ju.Node[mynodeval]

type nodetype1 struct {
	good bool
}

type nodetype2 struct {
	great bool
}

func (n nodetype1) isNodeVal() {}

func (n nodetype2) isNodeVal() {}

func (n *nodetype1) say() string {
	if n.good == true {
		return "good!"
	} else {
		return "boo"
	}
}

func (n *nodetype2) respond() int {
	if n.great {
		return 1
	} else {
		return 0
	}
}

func main() {

	v1 := nodetype1{good: true}

	v2 := nodetype2{great: false}

	n1 := new(mynode)

	n1.Val = v1

	n2 := new(mynode)

	n2.Val = v2

	n1.AddChild(n2)

	e := n1.Child(0)

	_, ok := e.Val.(nodetype2)

	if !ok {
		fmt.Println("FAIL")
		return
	}

	fmt.Println(e.Val.(nodetype2))
}
