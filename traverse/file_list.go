package traverse

import "sync"

type FileListRes struct {
	List  []string
	Map   map[string]struct{}
	ch    chan string
	mutex sync.Mutex
	*option
	over chan struct{}
}

func (s *FileListRes) callbackForGetFileList(item *TraverseItem) {

	if s.option.TargetExt != "" {
		if item.Ext == s.option.TargetExt {
			s.ch <- item.Path
		}
	}
}

func (s *FileListRes) processCh() {
	for {
		select {
		case p, ok := <-s.ch:
			if ok {
				s.mutex.Lock()
				s.List = append(s.List, p)
				s.Map[p] = struct{}{}
				s.mutex.Unlock()
			} else {
				s.over <- struct{}{}
				return
			}
		}
	}
}

func GetFileList(dir string, opts ...Option) (*FileListRes, error) {
	op := &option{}

	for _, o := range opts {
		o(op)
	}

	res := &FileListRes{
		Map:  make(map[string]struct{}),
		List: make([]string, 0),
		ch:   make(chan string, 512),
		over: make(chan struct{}),
	}
	res.option = op

	go res.processCh()

	t := NewDirTraverse(dir, res.callbackForGetFileList)
	err := t.Handle(opts...)
	if err != nil {
		close(res.ch)
		<-res.over
		return res, err
	}
	t.WorkSheet.Wait()
	close(res.ch)
	<-res.over
	return res, nil
}
