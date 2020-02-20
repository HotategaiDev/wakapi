package models

type Config struct {
	Env             string
	Port            int
	Addr            string
	DbHost          string
	DbPort          uint
	DbUser          string
	DbPassword      string
	DbName          string
	DbDialect       string
	DbMaxConn       uint
	CustomLanguages map[string]string
}

func (c *Config) IsDev() bool {
	return c.Env == "dev"
}
