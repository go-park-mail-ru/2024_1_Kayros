package dto

import "2024_1_kayros/internal/entity"

type Statistic struct {
	QuestionId   uint64 `json:"question_id" valid:"id"`
	QuestionName string `json:"question_name" valid:"question_name_domain"`
	Count        uint32 `json:"count" valid:"quiz_count_domain"`
	NPS          int8   `json:"nps" valid:"-"`
	CSAT         int8   `json:"csat" valid:"-"`
}

type StatisticArray struct {
	Payload []*Statistic `json:"payload" valid:"-"`
}

func NewDtoStatistic(statArray []*entity.Statistic) []*Statistic {
	statsDtoArray := []*Statistic{}
	for _, stat := range statArray {
		statsDtoArray = append(statsDtoArray, &Statistic{
			QuestionId:   stat.QuestionId,
			QuestionName: stat.QuestionName,
			Count:        stat.Count,
			NPS:          stat.NPS,
			CSAT:         stat.CSAT,
		})
	}
	return statsDtoArray
}
