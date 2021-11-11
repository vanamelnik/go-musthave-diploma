package psql

import "errors"

const defaultDSN = "postgres://root:qwe123@localhost:5432/gophermart?sslmode=disable"

type Config struct {
	DSN string `mapstructure:"uri"`
}

var defaultConfig Config = Config{
	DSN: defaultDSN,
}

func (c Config) Validate() error {
	if c.DSN == "" {
		return errors.New("storage: empty DSN")
	}

	return nil
}
