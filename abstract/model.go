package abstract

type Model struct {
}

func (this *Model) Options () map[string]interface{} {
	return map[string]interface{}{}
}

/*func (this *Model) TableName () string {
	return this.
}*/

func (this *Model) Source () string {
	return "default"
}

func (this *Model) Table () string {
	return ""
}