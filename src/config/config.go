package config

type (
	Config struct {
		Env              string           `mapstructure:"ENV"`
		Host             Host             `mapstructure:",squash"`
		PostgresDBConfig PostgresDBConfig `mapstructure:",squash"`
		MongoDBConfig    MongoDBConfig    `mapstructure:",squash"`
	}

	Host struct {
		Address      string `mapstructure:"HOST_ADDRESS"`
		Port         string `mapstructure:"HOST_PORT"`
		WriteTimeout int    `mapstructure:"HOST_WRITE_TIMEOUT"`
		ReadTimeout  int    `mapstructure:"HOST_READ_TIMEOUT"`
		IdleTimeout  int    `mapstructure:"HOST_IDLE_TIMEOUT"`
	}

	PostgresDBConfig struct {
		Host     string `mapstructure:"POSTGRES_DB_HOST"`
		User     string `mapstructure:"POSTGRES_DB_USER"`
		Password string `mapstructure:"POSTGRES_DB_PASSWORD"`
		Name     string `mapstructure:"POSTGRES_DB_NAME"`
		Port     string `mapstructure:"POSTGRES_DB_PORT"`
		SSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`
		Timezone string `mapstructure:"POSTGRES_TZ"`
	}

	MongoDBConfig struct {
		ConnectionString string `mapstructure:"MONGODB_URL"`
		DatabaseName     string `mapstructure:"MONGODB_DB_NAME"`
	}
)
