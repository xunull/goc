package traverse_v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/xunull/goc/traverse"
)

// buildBenchTree creates a tree with the given branching factor and depth.
// At each non-leaf level there are `dirs` subdirectories and `files` files.
func buildBenchTree(tb testing.TB, root string, depth, dirs, files int) {
	if depth == 0 {
		return
	}
	for i := 0; i < files; i++ {
		if err := os.WriteFile(filepath.Join(root, fmt.Sprintf("f%d.txt", i)), nil, 0644); err != nil {
			tb.Fatal(err)
		}
	}
	for i := 0; i < dirs; i++ {
		d := filepath.Join(root, fmt.Sprintf("d%d", i))
		if err := os.Mkdir(d, 0755); err != nil {
			tb.Fatal(err)
		}
		buildBenchTree(tb, d, depth-1, dirs, files)
	}
}

// Wide-shallow tree: 10 branching, 4 deep → 11,110 dirs + 100,000 files.
func BenchmarkV2_Traverse(b *testing.B) {
	tmp := b.TempDir()
	buildBenchTree(b, tmp, 4, 10, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var n atomic.Int64
		trv := New(tmp, func(item *Item) { n.Add(1) })
		if err := trv.Run(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkV1_Traverse(b *testing.B) {
	tmp := b.TempDir()
	buildBenchTree(b, tmp, 4, 10, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var n atomic.Int64
		dir := traverse.NewDirTraverse(tmp, func(item *traverse.TraverseItem) {
			n.Add(1)
		})
		if err := dir.Handle(); err != nil {
			b.Fatal(err)
		}
		dir.WorkSheet.Wait()
		dir.Close()
	}
}
