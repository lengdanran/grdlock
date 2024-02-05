// Package lk 2024/2/5 10:22
package lk

import (
	"context"
	"fmt"
	"github.com/lengdanran/grdlock/options"
	"github.com/redis/go-redis/v9"
	"time"
)

// RenewalRDL 自动续期分布式锁watchdog，会携带一个唯一标识作为value，解锁通过执行lua脚本判断是否为加锁者
type RenewalRDL struct {
	DefaultRDL
}

func (lk *RenewalRDL) WatchDog(opt *options.LOp) {
	ticker := time.NewTicker(opt.ExpireTime / 2)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-opt.DogQuitChan:
			// quit
			fmt.Println("Dog quit...")
			return
		default:
			fmt.Println("Dog renewal the lock------")
			lk.renewLock(opt)
		}
	}
}

func (lk *RenewalRDL) renewLock(opt *options.LOp) {
	const LuaCheckAndExpireDistributionLock = `
	    local lockerKey = KEYS[1]
	    local targetToken = ARGV[1]
	    local duration = ARGV[2]
	    local getToken = redis.call('get',lockerKey)
	    if (not getToken or getToken ~= targetToken) then
		    return 0
	    else
		    return redis.call('expire',lockerKey,duration)
	    end`
	keys := []string{opt.Key}
	eval := lk.Rds.Client.Eval(context.Background(), LuaCheckAndExpireDistributionLock, keys, opt.Value, opt.ExpireTime)
	if eval.Err() != nil {
		fmt.Println(eval.Err())
	}
	result, err := eval.Result()
	if err != nil {
		fmt.Println(err)
	}
	if ret, _ := result.(int64); ret != 1 {
		fmt.Println("can not expire lock without ownership of lock")
	}
}

func (lk *RenewalRDL) Lock(opt *options.LOp) error {
	err := lk.Rds.SetArgs(opt.Key, opt.Value, redis.SetArgs{Mode: "NX", TTL: opt.ExpireTime})
	if err != nil {
		return err
	}
	go lk.WatchDog(opt)
	return nil
}

func (lk *RenewalRDL) UnLock(opt *options.LOp) error {
	const UnLockScript = `
	    local lockerKey = KEYS[1]
	    local targetToken = ARGV[1]
	    local getToken = redis.call('get',lockerKey)
	    if (not getToken or getToken ~= targetToken) then
		    return 0
	    else
		    return redis.call('del',lockerKey)
	    end`
	keys := []string{opt.Key}
	eval := lk.Rds.Client.Eval(context.Background(), UnLockScript, keys, opt.Value)
	if eval.Err() != nil {
		fmt.Println(eval.Err())
		return eval.Err()
	}
	result, err := eval.Result()
	if err != nil {
		fmt.Println(err)
	}
	if ret, _ := result.(int64); ret != 1 {
		fmt.Println("can not del lock")
	}
	return err
}
