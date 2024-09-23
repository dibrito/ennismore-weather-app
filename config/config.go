package config

type ServiceConfig struct {
	APIConfig           APIConfig              `yaml:"api"`
	OpenstreetmapConfig OpenstreetmapAPIConfig `yaml:"openstreetmap"`
	WeatherConfig       WeatherAPIConfig       `yaml:"weather"`
}

type APIConfig struct {
	Port            int `yaml:"port"`
	ShutdownTimeout int `yaml:"shutdowntimeout"`
}

type OpenstreetmapAPIConfig struct {
	URL     string `yaml:"host"`
	Timeout int    `yaml:"timeout"`
}

type WeatherAPIConfig struct {
	URL     string `yaml:"host"`
	Timeout int    `yaml:"timeout"`
}
