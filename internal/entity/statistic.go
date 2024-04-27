package entity

type Statistic struct {
	QuestionId   uint64 `json:"question_id"`
	QuestionName string `json:"question_name"`
	Count        uint32 `json:"count"`
	NPS          int8   `json:"nps"`
	CSAT         int8   `json:"csat"`
}
