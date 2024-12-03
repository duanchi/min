package context

type Header map[string]string

func (this Header) Get(key string) string {
	return this[key]
}

func (this Header) Set(key, value string) {
	this[key] = value
}
