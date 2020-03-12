package form

type SaveClass struct {
	Id       int    `form:"id" json:"id"`
	ParentId int    `form:"parent_id" json:"parent_id"`
	Name     string `form:"name" json:"name"`
}

type ListClass struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Name  string `form:"name"`
}
