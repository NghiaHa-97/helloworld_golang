package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex  sync.Mutex
	amount int
)

//lock goroutine cho ddeesn khi mở khóa để không ảnh hưởng kết quả
func deposit(value int, wait *sync.WaitGroup) {
	mutex.Lock()
	fmt.Println("unlock deposit")
	time.Sleep(1 * time.Second)
	amount += value
	fmt.Println("deposit ", value)
	mutex.Unlock()
	defer wait.Done()
}

func withdraw(value int, wait *sync.WaitGroup) {
	mutex.Lock()
	fmt.Println("unlock withdraw")
	time.Sleep(7 * time.Second)
	amount -= value
	fmt.Println("withdraw ", value)
	mutex.Unlock()
	defer wait.Done()
}

func main() {
	var wait sync.WaitGroup
	amount = 1000
	wait.Add(2)
	go deposit(500, &wait)
	go withdraw(800, &wait)

	wait.Wait()
	fmt.Println(amount)
}
