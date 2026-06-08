package config

var conf *Config = nil

func UpdateSingletionConfig() {
	conf = Load()
}

func GetSingletonConfig() *Config {
	if conf == nil {
		UpdateSingletionConfig()
	}
	return conf
}
