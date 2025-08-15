package main

import (
	"fmt"
	"strings"

	"github.com/adamay909/logicTools/ju"
)

type mynodeval struct {
	content string
}

type mynode = ju.Node[mynodeval]

type mynode2 struct {
	mynode
	test string
}

func parse(s string) *mynode2 {

	prevnode := new(mynode2)
	n := new(mynode2)

	for i, c := range s {

		n = new(mynode2)

		n.test = "hello"

		n.Val.content = string(c)

		if i == 0 {
			prevnode = n
			continue
		}

		openAncestor(prevnode).AddChild(n)

		prevnode = n
	}

	return prevnode.Root()
}

func openAncestor(n *mynode2) *mynode2 {

	e := new(mynode2)

	for e = n; e != nil; e = e.Parent {

		switch e.Val.content {

		case "C", "K", "A":
			if len(e.Children()) < 2 {
				return e
			}

		case "N":
			if len(e.Children()) == 0 {
				return e
			}

		default:
			continue
		}
	}

	return nil
}

func render(n *mynode2) string {

	fmt.Println("rendering")
	var ingressFunc func(*mynode2)

	w := new(strings.Builder)

	ingressFunc = func(e *mynode2) {
		w.WriteString(e.Val.content)
	}

	donothing := func(e *mynode2) {
		return
	}

	n.Walk(ingressFunc, donothing, donothing)

	return w.String()
}

func isConditional(n *mynode) bool {
	return n.Val.content == "C"
}

func stringer(n *mynode) string {

	return "success: " + render(n)
}

func negateFunc(n any, v ...any) {

	m, ok := n.(*mynode)

	if !ok {
		panic("ju: improper use of function double")
	}

	p := new(mynode)

	p.Val.content = "N"

	p.AddChild(m)
}

func main() {

	ju.SetStringFunc[mynodeval](stringer)

	s := "CKCpqNqNp"

	var nd mynode2

	fmt.Println("processing", s)

	nd.ast = parse(s)

	nd.test = "hello"

	fmt.Println(nd.ast)

	fmt.Println(nd.test)

}
