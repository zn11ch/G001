package app

type Application struct {
	Dir   string `json:"dir"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Rules []struct {
		LogEveryMessage    bool `json:"logEveryMessage,omitempty"`
		SleepBeforeSending int  `json:"sleepBeforeSending,omitempty"`
	} `json:"rules"`
	Workers int `json:"workers"`
}

type Config struct {
	Debug        bool          `json:"debug"`
	Applications []Application `json:"applications"`
}

func NewConfig() *Config {

	c := Config{
		Debug:        false,
		Applications: []Application{},
	}
	return &c
}
