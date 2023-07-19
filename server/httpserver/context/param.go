package context

type Params struct {
	params map[string]string
}

func (this *Params) Get(key string, defaults ...string) string {
	if value, has := this.params[key]; has {
		return value
	} else {
		if len(defaults) > 0 {
			return defaults[0]
		} else {
			return ""
		}
	}
}

func (this *Params) GetAll() map[string]string {
	return this.params
}
