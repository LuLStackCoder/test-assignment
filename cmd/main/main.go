package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	counter2 "github.com/LuLStackCoder/test-assigment/pkg/counter"
	"github.com/LuLStackCoder/test-assigment/pkg/urlhandler"
)

const (
	errorMessage = "msg: error from count"
	outputPhrase = "Count for %s: %d"
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

	counter := counter2.NewCounter(httpClient, errorMessage, logger, outputPhrase)
	urlHandler := urlhandler.NewURLHandler(counter, rateLim)

	res, err := urlHandler.CountAllUrls(input, objectiveString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total:", res)
}
