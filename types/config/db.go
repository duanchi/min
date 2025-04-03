package config

type Db struct {
	Enabled      bool                `yaml:"enabled" default:"false"`
	Dsn          string              `yaml:"dsn"`
	MigrateSQL   string              `yaml:"migrate_sql"`
	Sources      map[string]DbConfig `yaml:"sources"`
	SelectEngine func(name string) (dsn string, err error)
}

type DbConfig struct {
	Dsn        string `yaml:"dsn"`
	MigrateSQL string `yaml:"migrate_sql"`
}
