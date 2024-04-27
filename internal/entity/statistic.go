package entity

type Statistic struct {
	QuestionId   uint16
	QuestionName string
	Count        uint32
	NPS          uint8
	CSAT         uint8
}
