// Package lk 2024/2/5 10:12
package lk

import (
	"github.com/lengdanran/grdlock/options"
	"github.com/redis/go-redis/v9"
)

// EXRDL 带过期时间的分布式锁实现 set nx ex/px
type EXRDL struct {
	DefaultRDL
}

func (lk *EXRDL) Lock(opt *options.LOp) error {
	err := lk.Rds.SetArgs(opt.Key, opt.Value, redis.SetArgs{Mode: "NX", TTL: opt.ExpireTime})
	if err != nil {
		return err
	}
	return nil
}
