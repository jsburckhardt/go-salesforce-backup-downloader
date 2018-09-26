package main

import (
	"fmt"
	"sync"
)

func worker(tasksCh <-chan string, wg *sync.WaitGroup, token string) {
	defer wg.Done()
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}

		fmt.Println("       Get-file-size(", task, ") using token:", token)
		fmt.Println("       download-file(", task, ") using token:", token)
		fmt.Println("Validating-file-size(", task, ")\n")

	}
}

func pool(wg *sync.WaitGroup, workers int, myvalues []string) {
	tasksCh := make(chan string)

	for i := 0; i < workers; i++ {
		go worker(tasksCh, wg, "hjkjdfad")
	}

	for _, v := range myvalues {
		tasksCh <- v
	}

	close(tasksCh)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	myvalues := []string{"url1", "url2", "url3", "url4", "url5", "url6", "url7", "url8", "url9"}
	go pool(&wg, 5, myvalues)
	wg.Wait()
}
