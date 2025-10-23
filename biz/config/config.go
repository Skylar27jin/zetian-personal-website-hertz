package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type SpecificConfig struct {
    Env          string `yaml:"env"`
    DB_DSN       string `yaml:"db_dsn"`
    Domain       string `yaml:"domain"`
}

type GeneralConfig struct {
	JWT_Secret_Key string `yaml:"jwt_secret_key"`
}

var specificCfg *SpecificConfig
var generalCfg *GeneralConfig

func InitConfig() {
    env := os.Getenv("ENV")
    if env == "" {
        env = "prod"
    }
    path := getPath(fmt.Sprintf("%s.yaml", env))
    data, err := os.ReadFile(path)
    if err != nil {
        panic(fmt.Sprintf("Failed to read config file %s: %v", path, err))
    }

    if err := yaml.Unmarshal(data, &specificCfg); err != nil {
        panic(fmt.Sprintf("Failed to parse YAML config: %v", err))
    }

    path = getPath("general_config.yaml")
    data, err = os.ReadFile(path)
    if err != nil {
        panic(fmt.Sprintf("Failed to read config file %s: %v", path, err))
    }
    if err := yaml.Unmarshal(data, &generalCfg); err != nil {
        panic(fmt.Sprintf("Failed to parse YAML config: %v", err))
    }

}

func GetSpecificConfig() SpecificConfig {
    return SpecificConfig{
        Env:    specificCfg.Env,
        DB_DSN: specificCfg.DB_DSN,
        Domain: specificCfg.Domain,
    }
}

func GetGeneralConfig() GeneralConfig {
    return GeneralConfig{
        JWT_Secret_Key: generalCfg.JWT_Secret_Key,
    }
}

func getPath(name string) string {
    cwd, _ := os.Getwd()
    return filepath.Join(cwd, "biz", "config", name)
}