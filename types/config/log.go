package config

type Log struct {
	Enabled bool `yaml:"enabled" default:"true"`
	Format struct {
		Type string `yaml:"type" default:"text"`
		Timestamp string `yaml:"timestamp" default:"2006-01-02 15:04:05.000"`
		Text FormatText `yaml:"text"`
		Json FormatJson `yaml:"json"`
	} `yaml:"format"`
	Timestamp bool `yaml:"timestamp" default:"true"`
	Output string `yaml:"output" default:"stdout://"`
	Level string `yaml:"level" default:"error"`
}

type FormatText struct {
	Colors bool `yaml:"colors" default:"true"`
	FullTimestamp bool `yaml:"full-timestamp" default:"true"`
}

type FormatJson struct {
	Pretty bool	`yaml:"pretty" default:"true"`
}
