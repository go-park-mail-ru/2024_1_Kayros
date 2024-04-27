package entity

type Statistic struct {
	QuestionId   uint16 `json:"question_id"`
	QuestionName string `json:"question_name"`
	Count        uint32 `json:"count"`
	NPS          uint8  `json:"nps"`
	CSAT         uint8  `json:"csat"`
}
