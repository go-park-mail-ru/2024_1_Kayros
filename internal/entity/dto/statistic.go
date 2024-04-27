package dto

type Statistic struct {
	QuestionId   uint16 `json:"question_id" valid:"id"`
	QuestionText string `json:"question_text" valid:"question_text_domain"`
	Count        uint32 `json:"count" valid:"quiz_count_domain"`
	Rating       uint8  `json:"rating" valid:"quiz_rating_domain"`
}
