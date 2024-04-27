package entity

type Statistic struct {
	QuestionId   uint16
	QuestionName string
	Count        uint32
	NPS          int8
	CSAT         int8
}
