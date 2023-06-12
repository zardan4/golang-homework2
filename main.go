package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	startTime := time.Now()
	const (
		usersCount   = 1000
		workersCount = 1000
	)

	jobs := make(chan int, usersCount)
	users := make(chan User, usersCount)

	wg := &sync.WaitGroup{}

	generateJobs(usersCount, jobs, wg)

	generateUsers(workersCount, jobs, users)

	saveUserInfo(workersCount, users, wg)

	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(workersCount int, users <-chan User, wg *sync.WaitGroup) {
	for i := 0; i < workersCount; i++ {
		go func() {
			for u := range users {
				fmt.Printf("WRITING FILE FOR UID %d\n", u.id)

				filename := fmt.Sprintf("users/uid%d.txt", u.id)
				file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}

				file.WriteString(u.getActivityInfo())
				time.Sleep(time.Second)

				wg.Done()
			}
		}()
	}
}

func generateJobs(count int, jobs chan<- int, wg *sync.WaitGroup) {
	wg.Add(count)

	for i := 0; i < count; i++ {
		jobs <- i
	}
}

func generateUsers(countWorkers int, jobs <-chan int, users chan<- User) {
	for i := 0; i < countWorkers; i++ {
		go func() {
			for i := range jobs {
				users <- User{
					id:    i + 1,
					email: fmt.Sprintf("user%d@company.com", i+1),
					logs:  generateLogs(rand.Intn(1000)),
				}
				fmt.Printf("generated user %d\n", i+1)
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
