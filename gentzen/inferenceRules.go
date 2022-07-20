package gentzen

import (
	"sort"
	"strings"
)

type sequent struct {
	datum, succedent string
}

var strictCheck bool

func datumReduce(d string) (out string) {

	d = strings.ReplaceAll(d, " ", "")
	data := strings.Split(d, ",")
	if len(data) == 0 {
		return d
	}

	sort.Strings(data)

	var ndata []string

	ndata = append(ndata, data[0])

	for i := 1; i < len(data); i++ {
		if ndata[len(ndata)-1] != data[i] {
			ndata = append(ndata, data[i])
		}
	}
	for i := range ndata {
		out = out + ndata[i] + ","
	}

	return strings.TrimRight(out, ",")
}

func equal(d1, d2 string) bool {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	if len(data1) != len(data2) {
		return false
	}

	sort.Strings(data1)
	sort.Strings(data2)

	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[i] {
			return false
		}
	}

	return true
}

//remove d2 from d1
func rm(d1, d2 string) string {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	var data3 []string

	for i := range data1 {
		if !in(data1[i], data2) {
			data3 = append(data3, data1[i])
		}
	}
	var d3 string

	for _, j := range data3 {
		d3 = d3 + j + ","
	}
	d3 = strings.TrimRight(d3, ",")
	return d3
}

func add(d1, d2 string) string {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")
	if d1 == "" {
		return d2
	}

	if d2 == "" {
		return d1
	}

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	for _, i := range data2 {
		data1 = append(data1, i)
	}

	sort.Strings(data1)

	var d3 string

	for _, j := range data1 {
		d3 = d3 + j + ","
	}
	d3 = strings.TrimRight(d3, ",")
	return d3
}

func contains(d string, dg string) bool {

	dg = strings.ReplaceAll(dg, " ", "")
	data1 := strings.Split(dg, ",")

	return in(d, data1)

}

func in(s string, g []string) bool {
	if s == "" {
		return true
	}
	for _, i := range g {
		if s == i {
			return true
		}
	}
	return false
}

func getSubnodes(n *Node) []*Node {

	var gs func(n *Node, list []*Node) []*Node

	gs = func(n *Node, list []*Node) []*Node {

		list = append(list, n)

		if n.IsAtomic() {
			return list
		}

		if n.subnode1 != nil {
			list = gs(n.subnode1, list)
		}

		if n.subnode2 != nil {
			list = gs(n.subnode2, list)
		}

		return list
	}

	var list []*Node

	return gs(n, list)
}

//order nodes by depth
func reorderNodes(nodes []*Node) (out []*Node) {

	d := findMaxDepth(nodes)

	for i := 0; i <= d; i++ {

		for _, j := range nodes {
			if j.Generation() == i {
				out = append(out, j)
			}
		}
	}
	return out
}

func findMaxDepth(nodes []*Node) int {

	var ds []int

	for _, n := range nodes {
		ds = append(ds, n.Generation())
	}

	sort.Ints(ds)

	return ds[len(ds)-1]
}

//check if s2 is instance of s1
func sameStructure(s1, s2 string) bool {

	ns1 := getSubnodes(Parse(s1))
	ns2 := getSubnodes(Parse(s2))

	ns1 = reorderNodes(ns1)
	ns2 = reorderNodes(ns2)

	if len(ns2) < len(ns1) {
		return false
	}

	for i := range ns1 {

		if ns1[i].IsAtomic() {
			old := ns1[i].Formula()
			repl := ns2[i].Formula()
			for _, n := range ns1 {
				if n.IsAtomic() && !n.HasFlag("c") {
					if n.Formula() == old {
						n.SetFormula(repl)
						n.SetFlag("c")
					}
				}
			}
		}
		continue

		if ns1[i].MainConnective() != ns2[i].MainConnective() {
			return false
		}
	}

	return printNodePolish(ns1[0]) == printNodePolish(ns2[0])

}
