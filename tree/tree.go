package tree

import (
	"fmt"
	"github.com/xunull/goc/vision"
	"reflect"
	"strings"
)

type ParentMap struct {
	M map[string][]*TreeElem
}

func (pm *ParentMap) Append(p string, item *TreeElem) {
	if _, ok := pm.M[p]; !ok {
		pm.M[p] = make([]*TreeElem, 0)
	}
	pm.M[p] = append(pm.M[p], item)
}

type LevelParentMap struct {
	M map[int]string
}

// ---------------------------------------------------------------------------------------------------------------------

type TreeElem struct {
	Data  interface{} // dir no data
	Next  []*TreeElem
	Leaf  bool
	Path  string
	Name  string
	Level int
	lpm   *LevelParentMap // no use
	pm    *ParentMap      // no use
}

func (t *TreeElem) IsRoot() bool {
	return t.Level == 0
}

func NewRooElem() *TreeElem {
	t := &TreeElem{Path: "", Level: 0}
	t.lpm = &LevelParentMap{
		M: make(map[int]string),
	}
	t.pm = &ParentMap{
		M: make(map[string][]*TreeElem),
	}
	return t
}

func listTreeElemMapData(target *TreeElem, v *vision.MindMapItem) {
	cur := vision.MindMapItem{
		Id:       target.Path,
		Children: make([]*vision.MindMapItem, 0),
	}
	v.Children = append(v.Children, &cur)
	for _, item := range target.Next {
		listTreeElemMapData(item, &cur)
	}
}

func (t *TreeElem) GetMindMapData() *vision.MindMapItem {
	// if Id use t.Name, may many some id in G6,G6 doesn't work
	root := vision.MindMapItem{
		Id: t.Path,
	}
	for _, item := range t.Next {
		listTreeElemMapData(item, &root)
	}
	return &root
}

func (t *TreeElem) PathString() string {
	return t.Name
}

func (t *TreeElem) Iter() []TreeAble {
	temp := make([]TreeAble, 0, len(t.Next))
	for _, item := range t.Next {
		temp = append(temp, item)
	}
	return temp
}

func NewElement(data interface{}) *TreeElem {
	ele := TreeElem{
		Data: data,
		Leaf: false,
	}
	return &ele
}

func appendNext(parent *TreeElem, kll [][]string, tm map[string]*TreeElem, root *TreeElem) {
	klength := len(kll)
	branchs := make(map[string]*TreeElem)
	names := make(map[string][][]string)
	for _, items := range kll {
		length := len(items)
		if length == 1 {
			// only full path
			return
		} else if length == 2 {
			// full path + leaf
			target := tm[items[0]]
			target.Name = items[1]
			target.Leaf = true
			target.Path = fmt.Sprintf("%s/%s", parent.Path, target.Name)
			target.Level = parent.Level + 1
			parent.Next = append(parent.Next, target)

			root.pm.Append(parent.Path, target)

		} else {
			// many

			name := items[1]
			if _, ok := names[name]; !ok {
				names[name] = make([][]string, 0, klength)
				branch := &TreeElem{
					Name:  name,
					Level: parent.Level + 1,
				}
				branchs[name] = branch
			}
			n := append([]string{items[0]}, items[2:]...)

			names[name] = append(names[name], n)
		}
	}

	for name, branch := range branchs {
		branch.Path = fmt.Sprintf("%s/%s", parent.Path, branch.Name)
		appendNext(branch, names[name], tm, root)
		parent.Next = append(parent.Next, branch)
	}

}

// ---------------------------------------------------------------------------------------------------------------------

func PlantPathTree(data []interface{}, key string) *TreeElem {
	eleList := make([]*TreeElem, 0, len(data))

	tm := make(map[string]*TreeElem)
	kll := make([][]string, 0, len(data))

	for _, item := range data {
		t := reflect.TypeOf(item)
		v := reflect.ValueOf(item)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			v = v.Elem()
		}

		if _, found := t.FieldByName(key); found {

			s := v.FieldByName(key).String()
			kl := strings.Split(s, "/")
			kl = append([]string{s}, kl...)
			kll = append(kll, kl)
			e := NewElement(item)
			eleList = append(eleList, e)
			tm[s] = e
		}
	}

	root := NewRooElem()
	appendNext(root, kll, tm, root)
	return root
}
