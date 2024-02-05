// Package options 2024/2/2 17:48
package options

import "time"

type LOp struct {
	Mode          string
	ExpireTime    time.Duration
	Key           string
	Value         string
	RetryTimes    int
	Retry         bool
	RetryInterval time.Duration
	DogQuitChan   chan bool
}

func DefaultLOp() *LOp {
	return &LOp{
		Mode:          "Simple",
		ExpireTime:    0,
		Key:           "default-lock-key",
		Value:         "",
		RetryTimes:    3,
		Retry:         true,
		RetryInterval: time.Millisecond * 500,
	}
}

func RenewalLockOpt(key, token string, exp time.Duration) *LOp {
	opts := DefaultLOp()
	opts.Mode = "Renewal"
	opts.Key = key
	opts.ExpireTime = exp
	opts.Value = token
	opts.DogQuitChan = make(chan bool, 1)
	return opts
}
