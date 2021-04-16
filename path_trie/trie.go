package path_trie

import (
	"strings"
)

type Trie struct {
	Separator string
	children  map[string]*Trie
	isEnd     bool
}

func NewTrie(separator string) *Trie {
	t := &Trie{
		Separator: separator,
		children:  make(map[string]*Trie),
	}
	return t
}

func (t *Trie) Insert(target string) {
	tl := strings.Split(target, "")
	node := t
	for _, item := range tl {
		if _, ok := node.children[item]; !ok {
			node.children[item] = NewTrie(t.Separator)
		}
		node = node.children[item]
	}
	node.isEnd = true
}

func (t *Trie) SearchPrefix(prefix string) *Trie {
	node := t
	pl := strings.Split(prefix, t.Separator)
	for _, item := range pl {
		if _, ok := node.children[item]; !ok {
			return nil
		}
		node = node.children[item]
	}
	return node
}

func (t *Trie) Search(target string) bool {
	node := t.SearchPrefix(target)
	return node != nil && node.isEnd
}
