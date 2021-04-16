package tree

import "fmt"

type TreeAble interface {
	Iter() []TreeAble
	PathString() string
}

type treeItResult struct {
	Data   []TreeAble
	Result []string
}

func (t *treeItResult) getPre(last bool, level int) string {
	if level == 0 {
		return ""
	}
	if last {
		return "    "
	} else {
		return "|   "
	}
}

func (t *treeItResult) processFull(ta TreeAble, pre string, last bool, level int) {
	var head string
	if last {
		head = "└── "
	} else {
		head = "├── "
	}
	if level == 0 {
		head = ""
	}

	full := fmt.Sprintf("%s%s%s", pre, head, ta.PathString())
	t.Result = append(t.Result, full)
}

func (t *treeItResult) process(ta []TreeAble, pre string, level int) {
	length := len(ta)
	for i, item := range ta {
		cur := t.getPre(i == length-1, level)
		if item.Iter() == nil {
			t.processFull(item, pre, i == length-1, level)
		} else {
			t.processFull(item, pre, i == length-1, level)
			t.process(item.Iter(), pre+cur, level+1)
		}
	}

}

func TreeIt(ta []TreeAble) []string {
	tir := treeItResult{
		Data: ta,
	}
	tir.process(tir.Data, "", 0)
	return tir.Result
}
