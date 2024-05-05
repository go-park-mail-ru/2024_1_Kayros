package entity

type Restaurant struct {
	Id               uint64
	Name             string
	ShortDescription string
	LongDescription  string
	Address          string
	ImgUrl           string
	Rating           float64
	CommentCount     uint32
}
