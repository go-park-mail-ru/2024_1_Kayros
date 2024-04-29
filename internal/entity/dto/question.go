package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

type QuestionInput struct {
	Id     uint64 `json:"id"`
	Rating uint8  `json:"rating"`
}

func (d *QuestionInput) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type Question struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	FocusId   string `json:"focus_id,omitempty"`
	ParamType string `json:"param_type"`
}

func (d *Question) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func QuestionReturn(qArray []*entity.Question) []Question {
	arr := []Question{}
	for _, q := range qArray {
		qDTO := Question{
			Id:        q.Id,
			Name:      q.Name,
			ParamType: q.ParamType,
			FocusId:   q.FocusId,
			Url:       q.Url,
		}
		arr = append(arr, qDTO)
	}
	return arr
}
