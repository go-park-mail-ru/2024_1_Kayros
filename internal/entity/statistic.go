package entity

type Statistic struct {
	QuestionId    uint16  `json:"question_id"`
	QuestionName  string  `json:"question_name"`
	Count         uint32  `json:"count"`
	AverageRating float32 `json:"average_rating"`
}
