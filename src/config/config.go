package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	Rpc      string `yaml:"rpc"`
}

func Load() (*Config, *ConfigFlags, error) {
	flags := GetConfigFlags()
	conf, err := ConfigFromYaml(flags.ConfigPath)
	return conf, flags, err
}

func ConfigFromYaml(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

type ConfigFlags struct {
	ConfigPath string
}

func GetConfigFlags() *ConfigFlags {
	o := &ConfigFlags{}
	flag.StringVar(&o.ConfigPath, "config", "config.yml", "path to config file")
	flag.Parse()
	return o
}
