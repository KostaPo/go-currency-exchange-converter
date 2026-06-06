package currency

type Currency struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"name"`
	Sign     string `json:"sign"`
}
