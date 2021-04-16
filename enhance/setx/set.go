package setx

type IntSet struct {
	m map[int]struct{}
}

func NewIntSet() *IntSet {
	return &IntSet{
		m: make(map[int]struct{}),
	}
}

func (in *IntSet) Add(target int) {
	in.m[target] = struct{}{}
}

func (in *IntSet) GetList() []int {
	res := make([]int, 0, len(in.m))
	for k, _ := range in.m {
		res = append(res, k)
	}
	return res
}
