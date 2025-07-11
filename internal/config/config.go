package config

type Config struct {
	// Temporary user shown by default - until we write a login page
	DefaultAppUsername string
	DBHost             string
	DBPort             int
	DBName             string
	DBUsername         string
	DBPassword         string
}

func LocalDev() Config {
	return Config{
		DefaultAppUsername: "salsajimmy",
		DBHost:             "localhost",
		DBPort:             5432,
		DBName:             "bite_tracker",
		DBUsername:         "biteapp",
		DBPassword:         "superburrito",
	}
}
