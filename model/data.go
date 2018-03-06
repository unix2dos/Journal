package model

type Data struct {
	Ret  int                    `json:"ret, string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func NewData() *Data {
	d := &Data{}
	d.Ret = Success
	d.Msg = ""
	d.Data = make(map[string]interface{})
	return d
}
