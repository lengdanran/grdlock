## grdlock

grdlock是一个简易的基于redis实现的分布式锁，可以根据不同场景使用不同的锁，提供默认实现，开箱即用。

### useage

```shell
go get github.com/lengdanran/grdlock
```

根据opt配置不同，使用统一的上锁/解锁入口即可使用不同类型的锁。

示例代码

```go
package main

import (
	"fmt"
	"github.com/lengdanran/grdlock/lk"
	"github.com/lengdanran/grdlock/options"
	"strconv"
	"time"
)

func lockTest(name string, ch chan int) {
	// RenewalLockOpt 为带有看门🐶自动锁续期机制的锁，默认会以指定过期时间的一半进行循环续期。
	// 当锁key被用户删除之后或者用户释放锁之后，看门🐶自动退出，解锁操作执行lua脚本，保证上锁者进行解锁。
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
		go lockTest("Thread-"+strconv.Itoa(i)+" ", ch)
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

```