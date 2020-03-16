package form

type ListUser struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Mobile string `form:"mobile"`
	Name   string `form:"name"`
}

type SaveUser struct {
	Id     int    `json:"id"`
	Mobile string `json:"mobile"`
	Passwd string `json:"passwd"`
	Name   string `json:"name"`
}
