package currency

type Currency struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
	Sign     string `json:"sign"`
}
