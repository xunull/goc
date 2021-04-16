package tree

func (t *TreeElem) traverse() {

}

type TraverseResultMap struct {
	M map[*TreeElem]interface{}
}

func (t *TreeElem) traverseAndHandleResult(tf func(data interface{}) interface{},
	hf func(data ...interface{}) interface{}, trm *TraverseResultMap) {

	resList := make([]interface{}, 0)

	for _, child := range t.Next {
		if child.Leaf {
			res := tf(child.Data)
			resList = append(resList, res)
		} else {
			child.traverseAndHandleResult(tf, hf, trm)
			res := trm.M[child]
			resList = append(resList, res)
		}
	}
	trm.M[t] = hf(resList...)
}

func (t *TreeElem) TraverseAndHandleResult(tf func(data interface{}) interface{},
	hf func(data ...interface{}) interface{},
	opts ...TreeOption) *TraverseResultMap {
	//option := getTreeOption(opts...)

	trm := &TraverseResultMap{
		M: make(map[*TreeElem]interface{}),
	}

	t.traverseAndHandleResult(tf, hf, trm)

	return trm
}
