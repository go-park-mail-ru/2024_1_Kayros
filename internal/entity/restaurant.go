package entity

type Restaurant struct {
	Id               uint64
	Name             string
	ShortDescription string
	LongDescription  string
	Address          string
	ImgUrl           string
	Rating           uint32
	CommentCount     uint32
}
