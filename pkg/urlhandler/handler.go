package urlhandler

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

// Counter interface required for URLHandler
type Counter interface {
	Count(url string, objectiveString string, totalValue *int64, outputChannel chan string)
}

// URLHandler implements URLHandler
type URLHandler struct {
	counter Counter
	maxCap  int
}

// CountAllUrls count the objective string from all pages read from reader
func (h *URLHandler) CountAllUrls(reader io.Reader, objectiveString string) (count int64, err error) {
	wgWorker := &sync.WaitGroup{}
	wgPrinter := &sync.WaitGroup{}

	scanner := bufio.NewScanner(reader)
	worker := 0

	workerChan := make(chan string, h.maxCap)
	outputChan := make(chan string, h.maxCap)

	wgPrinter.Add(1)

	go func(outputChan chan string) {
		defer wgPrinter.Done()
		for out := range outputChan {
			fmt.Println(out)
		}
	}(outputChan)

	for scanner.Scan() {
		if worker < h.maxCap {
			worker++
			go func(workerChan chan string) {
				for task := range workerChan {
					wgWorker.Add(1)
					h.counter.Count(task, objectiveString, &count, outputChan)
					wgWorker.Done()
				}
			}(workerChan)
		}

		workerChan <- scanner.Text()
	}

	close(workerChan)

	err = scanner.Err()
	if err != nil {
		return count, err
	}

	wgWorker.Wait()

	close(outputChan)

	wgPrinter.Wait()

	return
}

// NewURLHandler creates new URLHandler interface
func NewURLHandler(counter Counter, maxCap int) *URLHandler {
	return &URLHandler{
		maxCap:  maxCap,
		counter: counter,
	}
}
