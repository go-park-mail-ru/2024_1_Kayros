package entity

type Comment struct {
	Id        uint64
	UserId    uint64
	UserName  string
	UserImage string
	RestId    uint64
	Text      string
	Rating    uint32
}
