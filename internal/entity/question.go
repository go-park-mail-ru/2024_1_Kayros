package entity

type Question struct {
	Id        uint16 `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"-"`
	FocusId   string `json:"focus_id"`
	ParamType string `json:"param_type"`
}
