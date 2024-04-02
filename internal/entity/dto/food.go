package dto

import "github.com/asaskevich/govalidator"

type Food struct {
	Id          uint64 `json:"id" valid:"-"`
	Name        string `json:"name" valid:"-"`
	Description string `json:"description" valid:"-"`
	Restaurant  uint64 `json:"restaurant" valid:"-"`
	ImgUrl      string `json:"img_url" valid:"-"`
	Weight      uint64 `json:"weight" valid:"-"`
	Price       uint64 `json:"price" valid:"-"`
}

func (d *Food) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
