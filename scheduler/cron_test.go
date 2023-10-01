package scheduler

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	longTimeout = time.Minute
)

func should(t *testing.T, condition bool, format string, msg ...any) {
	if !condition {
		t.Errorf(format, msg...)
	}
}

// TestOne Tests the periodic repitition of a job
func TestOne(t *testing.T) {
	x := 0
	target := 0
	f := func() { x++ }

	c := NewCron("Test1")
	c.AddJob(longTimeout, time.Millisecond*100, f, 1)

	c.RunJob(1)

	time.Sleep(time.Millisecond * 10)
	for i := 0; i < 4; i++ {
		time.Sleep(time.Millisecond * 100)
		target++

		should(t, x == target, "got %q, wanted %q", x, target)
	}

	c.StopJob(1)
}

// TestTwo tests the corner cases of registering jobs to the job pool
func TestTwo(t *testing.T) {
	x := 0
	f1 := func() { x++ }
	f2 := func() { x-- }

	c := NewCron("Test2")

	//////////////////////////////////////
	c.AddJob(longTimeout, time.Millisecond*100, f1, 1)
	c.AddJob(longTimeout, time.Millisecond*100, f2, 1) // f2 should replace f1 in the job pool

	c.RunJob(1)

	time.Sleep(time.Millisecond * 200)

	should(t, x < 0, "f2 should have replaced f1")
	//////////////////////////////////////
	c.AddJob(longTimeout, time.Millisecond*10, f1, 1) // f1 should not replace f2 since it is already running

	time.Sleep(time.Millisecond * 200)
	should(t, x < 0, "f1 should not have replaced f2")
	//////////////////////////////////////
	c.StopJob(1)
	c.AddJob(longTimeout, time.Millisecond*10, f1, 1) // f1 should replace f2
	c.RunJob(1)

	time.Sleep(time.Millisecond * 200)
	should(t, x > 0, "f1 should have replaced f2")

	c.StopJob(1)
}

// TestThree tests that the failure of a scheduled job should halt the execution of another
func TestThree(t *testing.T) {
	zero := 0
	x := 0
	f1 := func() { x++ }
	f2 := func() { x = x / zero; x = 0 }

	c := NewCron("Test3")

	c.AddJob(longTimeout, 100*time.Millisecond, f1, 1)
	c.AddJob(longTimeout, 10*time.Millisecond, f2, 2)

	c.RunAll()

	time.Sleep(200 * time.Millisecond)
	should(t, x > 0, "f2 failure should not have stopped the scheduling of f1")

	c.StopJob(1)
}

// TestFour attempts to run a job not existing in the job pool, the scheduler should not crash then
func TestFour(t *testing.T) {
	x := 0
	f := func() { x++ }

	c := NewCron("Test4")
	c.AddJob(longTimeout, 10*time.Millisecond, f, 1)
	c.RunJob(2)
	c.RunJob(1)

	time.Sleep(50 * time.Millisecond)
	should(t, x > 0, "attemping to run a non-existing job should not affect the scheduling of other jobs, x = %d", x)

	c.StopJob(2)
	c.StopJob(1)
	c.StopJob(1)

	time.Sleep(10 * time.Millisecond)
	lastX := x
	time.Sleep(50 * time.Millisecond)
	should(t, x == lastX, "attemping to stop a non-existing job should not affect the scheduling of other jobs, x = %d, lastX = %d", x, lastX)
}

// TestFive attempts to run an already running job, the scheduler should not crash then
func TestFive(t *testing.T) {
	x := 0
	f := func() { x++ }

	c := NewCron("Test5")
	c.AddJob(longTimeout, 10*time.Millisecond, f, 1)

	c.RunJob(1)
	c.RunJob(1)

	time.Sleep(50 * time.Millisecond)
	should(t, x > 0, "attempting to run f multiple times should not crash the scheduler")

	c.StopJob(1)
}

// TestSix tests logging to console VS logging to a file
func simulateCron(file ...string) {
	c := NewCron(file...)
	c.AddJob(longTimeout, 10*time.Millisecond, func() {}, 1)
	c.RunJob(1)
	c.RunJob(1)                            // To have a warning in the logs
	c.AddJob(longTimeout, 0, func() {}, 1) // To have an error in logs
	time.Sleep(20 * time.Millisecond)
	c.StopJob(1)
}
func validateLogs(t *testing.T, logs string, consoleLogs bool) {
	allShouldContain := []string{
		"CRON",
		"INFO",
		"WARN",
		"ERROR",
		"2023",
		"Registered job with id 1 successfully!!!",
		"Started Job with id 1",
	}
	consoleShouldContain := []string{
		"\033[1;33m",
		"\033[1;31m",
		"\033[0m",
	}
	containsAll := func(arr []string) {
		for _, str := range arr {
			should(t, strings.Contains(logs, str), "the logs should have contained the text %s", str)
		}
	}

	containsAll(allShouldContain)
	if consoleLogs {
		containsAll(consoleShouldContain)
	}
}
func TestSix(t *testing.T) {
	////////////////////////////////////// console logs
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	simulateCron() // logs to console

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	should(t, err == nil, "reading from the console output failed")

	validateLogs(t, buf.String(), true)

	////////////////////////////////////// file logs
	fileName := "Test6"
	filePath := "logs/" + fileName + ".log"
	os.Remove(filePath)
	simulateCron(fileName)

	fileLogs, err := os.ReadFile(filePath)
	should(t, err == nil, "error reading log file at path %s", filePath)
	validateLogs(t, string(fileLogs), false)
}

// TestSeven tests scheduling a task for negative period, the period should default to 0. It also tries timeout jobs.
func TestSeven(t *testing.T) {
	x := 0
	f := func() { x++ }

	c := NewCron("Test7")
	c.AddJob(longTimeout, -100, f, 1)
	c.RunJob(1)

	time.Sleep(time.Millisecond)
	should(t, x > 0, "negative period should not crash the schedule but rather default to 0s, x = %d", x)
	c.StopJob(1)

	// function takes longer than its period
	y := 0
	g := func() { y--; time.Sleep(100 * time.Millisecond) }
	c.AddJob(longTimeout, 10*time.Millisecond, g, 1)
	c.RunJob(1)

	time.Sleep(50 * time.Millisecond)
	should(t, y == -1, "the period should be the job execution time if it takes longer that the period assigned to it, x = %d", x)
	c.StopJob(1)

	// job execution that times out
	z := 1
	h := func() { z *= 2; time.Sleep(2 * time.Millisecond) }
	c.AddJob(time.Millisecond, 10*time.Millisecond, h, 1)
	c.RunJob(1)

	time.Sleep(50 * time.Millisecond)
	should(t, z == 2, "the job execution should time out after first execution, z = %d", z)
	c.StopJob(1)
}

// TestEight tests scheduling multiple jobs together
func TestEight(t *testing.T) {
	x := 0
	y := 0
	z := 4

	targetX := 0
	targetY := 0

	f1 := func() { x++ }
	f2 := func() { y-- }
	f3 := func() { z--; z /= z }

	c := NewCron("Test8")
	c.AddJob(longTimeout, 1*time.Second, f1, 1)
	c.AddJob(longTimeout, 2*time.Second, f2, 2)
	c.AddJob(longTimeout, 1*time.Second, f3, 3)

	c.RunAll()

	time.Sleep(100 * time.Millisecond)
	for itr := 0; itr < 5; itr++ {
		time.Sleep(2 * time.Second)

		targetX += 2
		targetY -= 1

		should(t, x == targetX, "value of x => %d expected to be %d", x, targetX)
		should(t, y == targetY, "value of y => %d expected to be %d", y, targetY)
	}

	c.StopAll()
}
