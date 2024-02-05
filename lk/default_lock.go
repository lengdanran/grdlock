// Package lk 2024/2/4 16:20
package lk

import (
	"github.com/lengdanran/grdlock/options"
	"github.com/lengdanran/grdlock/rds"
	"github.com/redis/go-redis/v9"
)

// DefaultRDL 简单分布式锁实现 set nx,不带有过期时间
type DefaultRDL struct {
	Rds *rds.Rds
}

func (lk *DefaultRDL) Lock(opt *options.LOp) error {
	err := lk.Rds.SetArgs(opt.Key, opt.Value, redis.SetArgs{Mode: "NX"})
	if err != nil {
		return err
	}
	return nil
}

func (lk *DefaultRDL) UnLock(opt *options.LOp) error {
	err := lk.Rds.Del(opt.Key)
	if err != nil {
		return err
	}
	return nil
}
