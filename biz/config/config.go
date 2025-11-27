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
    CookieSecure bool `yaml:"cookie_secure"`


    //S3 / CloudFront 相关配置
    AWSRegion string `yaml:"aws_region"` // 比如 "us-east-2"
    S3Bucket  string `yaml:"s3_bucket"`  // 比如 "project-talk-media"
    CDNDomain string `yaml:"cdn_domain"` // 比如 "cdn.skylar27.com",可留空
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
        CookieSecure: specificCfg.CookieSecure,
        AWSRegion: specificCfg.AWSRegion,
        S3Bucket:  specificCfg.S3Bucket,
        CDNDomain: specificCfg.CDNDomain,
    }
}

func GetGeneralConfig() GeneralConfig {
    return GeneralConfig{
        JWT_Secret_Key: generalCfg.JWT_Secret_Key,
    }
}

func getPath(name string) string {
    cwd, _ := os.Getwd()


    basePath := filepath.Join(cwd, "biz", "config", name)

    // 如果在测试中（cwd 含 biz/service/...），则回退两层
    if _, err := os.Stat(basePath); os.IsNotExist(err) {
        basePath = filepath.Join(cwd, "..", "..", "config", name)
    }

    return basePath
}
