package model

type Category struct {
	BaseModel
	Name       string     `json:"name"`        // 分类名
	CategoryID int64      `json:"category_id"` // 父ID 默认0为一级
	Categories []Category `json:"categories"`  // 子分类列表
}
