package abstract

type Rpc struct {
	Service

	PackageName string
	ApplicationName string
}

func (this *Rpc) GetPackageName () (name string) {
	return this.PackageName
}

func (this *Rpc) SetPackageName (name string) {
	this.PackageName = name
}

func (this *Rpc) GetApplicationName () (name string) {
	return this.ApplicationName
}

func (this *Rpc) SetApplicationName (name string) {
	this.ApplicationName = name
}