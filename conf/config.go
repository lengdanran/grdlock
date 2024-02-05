// Package conf 2024/2/4 10:29
package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Redis RedisConf `json:"redis"`
}

type RedisConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	DB   int    `json:"db"`
	Pwd  string `json:"pwd"`
}

func GetConf(confFile string) (*Config, error) {
	dataBytes, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Println("读取配置文件失败：", err)
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return nil, err
	}
	mp := make(map[string]any, 2)
	err = yaml.Unmarshal(dataBytes, mp)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return nil, err
	}
	return &config, nil
}
