// Package rds 2024/2/2 17:40
package rds

import (
	"context"
	"fmt"
	"github.com/lengdanran/grdlock/conf"
	"github.com/redis/go-redis/v9"
)

type Rds struct {
	RdsOptions redis.Options
	Client     *redis.Client
}

func NewRdsWithConf(confFile string) *Rds {
	config, err := conf.GetConf(confFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	rdsOption := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Pwd,
		DB:       config.Redis.DB,
	}
	return &Rds{
		RdsOptions: rdsOption,
		Client:     redis.NewClient(&rdsOption),
	}
}

func NewRds() *Rds {
	return NewRdsWithConf("conf/config.yaml")
}

func (rds *Rds) NewRdsClient() *redis.Client {
	return redis.NewClient(&rds.RdsOptions)
}

func (rds *Rds) SetArgs(key string, value interface{}, args redis.SetArgs) error {
	err := rds.Client.SetArgs(context.Background(), key, value, args).Err()
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}
	return nil
}

func (rds *Rds) Eval(script string, keys []string, args ...interface{}) error {
	err := rds.Client.Eval(context.Background(), script, keys, args).Err()
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}
	return nil
}

func (rds *Rds) Get(key string) (string, error) {
	keyRet := rds.Client.Get(context.Background(), key)
	if keyRet.Err() != nil {
		// fmt.Println(keyRet.Err().Error())
		return "", keyRet.Err()
	}
	return keyRet.String(), nil
}

func (rds *Rds) Del(key string) error {
	delRet := rds.Client.Del(context.Background(), key)
	if delRet.Err() != nil {
		// fmt.Println(delRet.Err().Error())
		return delRet.Err()
	}
	return nil
}
