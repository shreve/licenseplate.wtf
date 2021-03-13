package server

type config struct {
	Port string
}

func loadConfig() *config {
	return &config{
		Port: ":8081",
	}
}
