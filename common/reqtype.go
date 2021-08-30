package common

type ReqType int

const (
	None ReqType = iota
	Get
	Post
	Put
	Delete
)

func (this ReqType) String() string {
	switch this {
	case Get:
		return "Pget"
	case Post:
		return "Ppost"
	case Put:
		return "Pput"
	case Delete:
		return "Pdelete"
	}
	return "N/A"
}



