package lock

import (
	"example/pkg/redisstore"
	"fmt"
	"testing"
	"time"
)

func TestNewDistributedLock(t *testing.T) {
	fmt.Println("start")
	//调用redis
	redisstore.Init()
	//测试分布式锁
	go locks()
	go locks()
	go locks()

	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("end")
}

func locks() {
	locks := NewDistributedLock()
	var (
		sleepMillisecond int
	)
	defer locks.UnLock()
LoopLock:
	ok, err := locks.TryLock()
	//分布锁进行自旋
	if err != nil || !ok {
		sleepMillisecond += 100
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println(sleepMillisecond)
		if sleepMillisecond > 50000 {
			return
		} else {
			goto LoopLock
		}
	}

	fmt.Println(2)
}
