package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vanamelnik/go-musthave-diploma/service/gophermart"

	"github.com/hashicorp/go-multierror"
	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

var defaultConfig = Config{
	RunAddr:           ":8080",
	AccrualSystemAddr: "localhost:3000",
	Logger: LoggerConfig{
		Level:   "debug",
		Console: false,
	},
	DatabaseURI: "postgres://root:secret@localhost:5432/gophermart?sslmode=disable",
	Service: gophermart.Config{
		PasswordPepper: "secret",
		UpdateInterval: 2 * time.Second,
	},
}

type (
	// Config represents configuration for all services.
	Config struct {
		RunAddr           string `mapstructure:"run_address"`
		AccrualSystemAddr string `mapstructure:"accrual_system_address"`

		Logger      LoggerConfig
		DatabaseURI string `mapstructure:"database_uri"`
		Service     gophermart.Config
	}

	LoggerConfig struct {
		Level   string `mapstructure:"level"`
		Console bool   `mapstructure:"console"`
	}

	Option func(cfg *Config)
)

func (c Config) Validate() (retErr error) {
	if c.RunAddr == "" {
		retErr = multierror.Append(retErr, errors.New("missing run address"))
	}
	if c.Service.UpdateInterval <= 0 {
		retErr = multierror.Append(retErr, errors.New("update interval is zero or less"))
	}
	if c.AccrualSystemAddr == "" {
		retErr = multierror.Append(retErr, errors.New("accrual system address not set"))
	}

	return retErr
}

// LoadConfig sets up the configuration loaded from the file provided, environment variables
// and flags.
func LoadConfig(cfgFileName string) Config {
	setDefaultConfig()
	viper.SetConfigFile(cfgFileName)
	bindFlags()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	c := Config{}
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config: could not load %s: %v, use default config\n", cfgFileName, err)
	}
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("config: could not unmarshal config: %v\n", err)
	}

	return c
}

func bindFlags() {
	_ = flag.StringP("run_address", "a", defaultConfig.RunAddr, "service's run address and port")
	_ = flag.StringP("database_uri", "d", defaultConfig.DatabaseURI, "DSN string for database connection")
	_ = flag.StringP("accrual_system_address", "r", defaultConfig.AccrualSystemAddr, "service's run address and port")
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		fmt.Printf("binding flags: %v\n", err)
	}
}

func setDefaultConfig() {
	viper.SetDefault("run_address", defaultConfig.RunAddr)
	viper.SetDefault("accrual_system_address", defaultConfig.AccrualSystemAddr)
	viper.SetDefault("logger.level", defaultConfig.Logger.Level)
	viper.SetDefault("logger.console", defaultConfig.Logger.Console)
	viper.SetDefault("database_uri", defaultConfig.DatabaseURI)
	viper.SetDefault("service.password_pepper", defaultConfig.Service.PasswordPepper)
	viper.SetDefault("service.update_interval", defaultConfig.Service.UpdateInterval)
}
