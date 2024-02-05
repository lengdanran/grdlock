// Package main 2024/2/2 17:34
package main

import (
	"fmt"
	"github.com/lengdanran/grdlock/lk"
	"github.com/lengdanran/grdlock/options"
	"strconv"
	"time"
)

func lockTest1(name string, ch chan int) {
	opts := options.RenewalLockOpt("TestLock", name, time.Second*4)
	if lk.Lock(opts) {
		fmt.Println(name + "Get lock OK!!!")
		fmt.Println(name + "Doing somethings")
		time.Sleep(time.Second * 10)

		if lk.ReleaseLock(opts) {
			fmt.Println(name + "Release lock OK")
		} else {
			fmt.Println(name + "Release lock Failed")
		}
		ch <- 1 // 协程结束，发出信号
		return
	}
	fmt.Println(name + "Not get Lock")
	ch <- 1 // 协程结束，发出信号
}

func main() {
	ch := make(chan int)
	count := 5 // 协程个数
	for i := 0; i < count; i++ {
		go lockTest1("Thread-"+strconv.Itoa(i)+" ", ch)
	}
	for range ch {
		count--
		if count == 0 {
			close(ch)
		}
	}
	time.Sleep(time.Second * 5)
	fmt.Println("Done")
}
