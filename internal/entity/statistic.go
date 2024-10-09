package entity

type Statistic struct {
	QuestionId   uint64
	QuestionName string
	Count        uint32
	NPS          int8
	CSAT         int8
}
