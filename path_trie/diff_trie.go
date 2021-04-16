package path_trie

func DiffPath(left, right []string) ([]string, []string, []string) {
	lt := NewTrie("/")
	rt := NewTrie("/")
	for _, item := range left {
		lt.Insert(item)
	}
	for _, item := range right {
		rt.Insert(item)
	}
	lm := make([]string, 0, len(left))
	rm := make([]string, 0, len(right))
	same := make([]string, 0, len(left))
	for _, item := range left {
		if rt.Search(item) {
			same = append(same, item)
		} else {
			lm = append(lm, item)
		}
	}

	for _, item := range right {
		if !lt.Search(item) {
			rm = append(rm, item)
		}
	}
	return lm, same, rm
}
