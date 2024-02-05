// Package lk 2024/2/2 17:35
package lk

import (
	"fmt"
	"github.com/lengdanran/grdlock/options"
	"github.com/lengdanran/grdlock/rds"
	"os"
	"time"
)

type L interface {
	// Lock 上锁，若无err产生，代表上锁成功
	Lock(opt *options.LOp) error
	// UnLock 解锁，若无err产生，代表解锁成功
	UnLock(opt *options.LOp) error
}

var LockerRegistry map[string]L
var Locker *RDL
var rdsConf = "conf/config.yaml"

func init() {
	envRdsConf := os.Getenv("GRDLOCK_CONF")
	if envRdsConf != "" {
		rdsConf = envRdsConf
	}
	LockerRegistry = make(map[string]L)
	LockerRegistry["Default"] = &DefaultRDL{Rds: rds.NewRdsWithConf(rdsConf)}

	ex := &EXRDL{}
	ex.Rds = rds.NewRdsWithConf(rdsConf)
	LockerRegistry["EX"] = ex

	renew := &RenewalRDL{}
	renew.Rds = rds.NewRdsWithConf(rdsConf)
	LockerRegistry["Renewal"] = renew

	Locker = &RDL{}
}

func GetLocker(mode string) L {
	if lock, ok := LockerRegistry[mode]; ok {
		return lock
	}
	// 处理未知的 mode 值或选择默认实现
	return &DefaultRDL{Rds: rds.NewRds()}
}

// RDL 分布式锁实现
type RDL struct {
}

func (rdl *RDL) execute(opt *options.LOp, execType string) error {
	var err error
	if opt.Retry {
		retryTimes := opt.RetryTimes
		if retryTimes <= 0 {
			retryTimes = 1
		}
		for i := 0; i < retryTimes; i++ {
			if execType == "Lock" {
				err = GetLocker(opt.Mode).Lock(opt)
			} else {
				err = GetLocker(opt.Mode).UnLock(opt)
			}
			if err != nil {
				if i == retryTimes-1 {
					return err
				}
			} else {
				return nil
			}
			time.Sleep(opt.RetryInterval)
			fmt.Println("Retry " + execType)
		}
	} else {
		if execType == "Lock" {
			err = GetLocker(opt.Mode).Lock(opt)
		} else {
			err = GetLocker(opt.Mode).UnLock(opt)
		}
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	return err
}

func (rdl *RDL) Lock(opt *options.LOp) error {
	return rdl.execute(opt, "Lock")
}

func (rdl *RDL) UnLock(opt *options.LOp) error {
	return rdl.execute(opt, "UnLock")
}

func Lock(opts *options.LOp) bool {
	err := Locker.Lock(opts)
	if err != nil {
		return false
	}
	return true
}

func ReleaseLock(opts *options.LOp) bool {
	defer func() {
		if opts.DogQuitChan != nil {
			opts.DogQuitChan <- true
		}
	}()
	err := Locker.UnLock(opts)
	if err != nil {
		return false
	}
	if err != nil {
		return false
	}
	return true
}
