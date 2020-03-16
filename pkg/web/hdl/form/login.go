package form

type LoginRequest struct {
	Mobile string `form:"mobile"`
	Passwd string `form:"passwd"`
}
