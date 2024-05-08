package dto

type Statistic struct {
	QuestionId   uint16 `json:"question_id" valid:"id"`
	QuestionName string `json:"question_name" valid:"question_name_domain"`
	Count        uint32 `json:"-" valid:"quiz_count_domain"`
	NPS          uint8  `json:"nps" valid:"-"`
	CSAT         uint8  `json:"csat" valid:"-"`
}
