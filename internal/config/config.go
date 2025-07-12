package config

type Config struct {
	// Temporary user shown by default - until we write a login page
	DefaultAppUserId string
	DBHost           string
	DBPort           int
	DBName           string
	DBUsername       string
	DBPassword       string
}

func LocalDev() Config {
	return Config{
		DefaultAppUserId: "f41ad27a-881d-4f7f-a908-f16a26ce7b78",
		DBHost:           "localhost",
		DBPort:           5432,
		DBName:           "bite_tracker",
		DBUsername:       "biteapp",
		DBPassword:       "superburrito",
	}
}
