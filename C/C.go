package C

type JSON struct {
	Err    int         `json:"err"`
	Data   interface{} `json:"data"`
	Reason string      `json:"reason"`
}

func Pack(v interface{}) *JSON {
	return &JSON{
		Err:    0,
		Data:   v,
		Reason: "",
	}
}

func PackError(v interface{}, e error) *JSON {
	if e != nil {
		return &JSON{
			Err:    1,
			Data:   nil,
			Reason: e.Error(),
		}
	}
	return Pack(v)
}