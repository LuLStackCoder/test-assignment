package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	counter2 "github.com/LuLStackCoder/test-assigment/pkg/counter"
	"github.com/LuLStackCoder/test-assigment/pkg/urlhandler"
)

const (
	errorMessage = "msg: error from count"
	outputPhrase = "Count for %s: %d\n"
)

func main() {
	var (
		rateLim         int
		objectiveString string
	)

	flag.IntVar(&rateLim, "k", 5, "num of goroutines")
	flag.StringVar(&objectiveString, "q", "go", "objective string")

	flag.Parse()

	input := os.Stdin

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)
	quotaChannel := make(chan struct{}, rateLim)
	counter := counter2.NewCounter(httpClient, quotaChannel, errorMessage, logger, os.Stdout, outputPhrase)
	urlHandler := urlhandler.NewURLHandler(counter, &sync.WaitGroup{})
	fmt.Println("Total:", urlHandler.CountAllUrls(input, objectiveString))
}
