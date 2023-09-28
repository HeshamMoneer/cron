package main

import (
	"fmt"
	"time"

	cron "cron/scheduler"
)

func main() {
	c := cron.CronInit()
	c.AddFunction(time.Millisecond*1000, func() {
		fmt.Println(1)
	})
	c.AddFunction(time.Millisecond*2000, func() {
		fmt.Println(2)
	})
	c.AddFunction(time.Second*10, func() {
		fmt.Println(3)
	})

	c.RunAll()

	fmt.Println("-------------")
}
