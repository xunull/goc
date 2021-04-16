package vision

type MindMapItem struct {
	Id       string         `json:"id"`
	Children []*MindMapItem `json:"children"`
}

type MindMap interface {
	GetMindMapData() *MindMapItem
}
