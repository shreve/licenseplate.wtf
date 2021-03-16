package server

type config struct {
	Port         string
	CacheEnabled bool
}

func loadConfig() *config {
	return &config{
		Port:         ":8081",
		CacheEnabled: true,
	}
}
