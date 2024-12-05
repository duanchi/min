package context

type Header map[string]string

func (this Header) Get(key string) string {
	return this[key]
}

func (this Header) Set(key, value string) {
	this[key] = value
}

func (this Header) Del(key string) {
	delete(this, key)
}

func (this Header) Clone() Header {
	clone := make(Header, len(this))
	for k, v := range this {
		clone[k] = v
	}
	return clone
}
