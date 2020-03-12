package form

type LoginRequest struct {
	Account string `form:"account"`
	Passwd  string `form:"passwd"`
}
