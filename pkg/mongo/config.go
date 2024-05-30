package mongo

type Config struct {
	Host string
	Port uint16
	Name string
	User string
	Pass string
}

func (c Config) IsEmptyUser() bool {
	return c.User == ""
}
