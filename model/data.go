package model

type Data struct {
	Ret  int                    `json:"ret, string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func NewData() *Data {
	d := &Data{}
	d.Ret = Success
	GetDataMsg(d, d.Ret)
	d.Data = make(map[string]interface{})
	return d
}
