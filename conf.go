package textra

type Config struct {
	BaseURL    string
	UserName   string
	API_KEY    string
	API_SECRET string
}

func NewConfig(username string, api_key string, api_secret string) *Config {
	return &Config{
		BaseURL:    BaseURL,
		UserName:   username,
		API_KEY:    api_key,
		API_SECRET: api_secret,
	}
}
