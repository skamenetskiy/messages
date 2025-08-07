package database

type Config struct {
	Shards []struct {
		ID       uint16 `yaml:"id"`
		Writable bool   `yaml:"writable"`
		DSN      string `yaml:"dsn"`
	} `yaml:"shards"`
}
