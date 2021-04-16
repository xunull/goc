package file_utils

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type FileTask struct {
	wg    *sync.WaitGroup
	Name  string
	Path  string
	IsDir bool
}

func TaskConsume(i interface{}) {
	task := i.(*FileTask)
	defer task.wg.Done()
	if task.IsDir {
		err := os.MkdirAll(task.Path, fs.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("target: %s created", task.Path)
		}
	} else {
		err := ioutil.WriteFile(task.Path,
			[]byte(gofakeit.Paragraph(5, 10, 200, "\n")),
			fs.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("target: %s created", task.Path)
		}
	}
}

func makeTempFiles(dirname string, width, depth int,
	p1Tasks, p2Tasks []*FileTask,
	wg1, wg2 *sync.WaitGroup) {
	if depth < 0 {
		return
	}
	fmt.Printf("%v %v %v \n", dirname, width, depth)
	for i := 0; i < width; i++ {
		n := gofakeit.LetterN(10)
		next := path.Join(dirname, n)
		wg1.Add(1)

		p1Tasks = append(p1Tasks, &FileTask{
			IsDir: true,
			Name:  n,
			Path:  path.Join(dirname, n),
			wg:    wg1,
		})

		makeTempFiles(next, width, depth-1, p1Tasks, p2Tasks, wg1, wg2)
	}
	for i := 0; i < width; i++ {
		n := gofakeit.LetterN(12)

		wg2.Add(1)
		p2Tasks = append(p2Tasks, &FileTask{
			IsDir: false,
			Name:  n,
			Path:  path.Join(dirname, n),
			wg:    wg2,
		})
	}
}

// width * depth
//func MakeTempFiles(dirname string, width, depth int) error {
//	if depth <= 0 {
//		return nil
//	}
//
//	p1, err := routine_pool.NewPool(32, TaskConsume)
//	defer p1.Release()
//
//	if err != nil {
//		return err
//	}
//	var wg1 sync.WaitGroup
//	var wg2 sync.WaitGroup
//	p1Tasks := make([]*FileTask, 0)
//	p2Tasks := make([]*FileTask, 0)
//	makeTempFiles(dirname, width, depth, p1Tasks, p2Tasks, &wg1, &wg2)
//	for _, task := range p1Tasks {
//		err := p1.Submit(task)
//		if err != nil {
//			return err
//		}
//	}
//	wg1.Wait()
//	for _, task := range p2Tasks {
//		err := p1.Submit(task)
//		if err != nil {
//			return err
//		}
//	}
//	wg2.Wait()
//	return nil
//}
