package dto

import "github.com/asaskevich/govalidator"

type Question struct {
	Id     uint64 `json:"id" valid:"-"`
	Rating uint32 `json:"rating" valid:"-"`
}

func (d *Question) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
