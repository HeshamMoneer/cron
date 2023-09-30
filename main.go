package main

import (
	"fmt"
	"time"

	cron "cron/scheduler"
)

func main() {
	c := cron.NewCron()
	c.AddJob(time.Millisecond, time.Millisecond*1000, func() {
		fmt.Println(1)
	}, 1)
	c.AddJob(time.Millisecond, time.Millisecond*2000, func() {
		fmt.Println(2)
	}, 2)
	c.AddJob(time.Millisecond, time.Second*10, func() {
		fmt.Println(3)
	}, 2)

	c.RunAll()

	c.WaitJobs()

	fmt.Println("-------------")
}
