package counter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)

// Counter interface represent the objective string counter on the single page
type Counter interface {
	Count(url string, objectiveString string, totalValue *int64, outputChannel chan string)
}

type counter struct {
	httpClient   *http.Client
	errorMessage string
	logger       *log.Logger
	outputPhrase string
}

// Count objective string on the single page
func (c *counter) Count(url string, objectiveString string, totalValue *int64, outputChan chan string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Println(c.errorMessage, err)
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Println(c.errorMessage, err)
		return
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Println(c.errorMessage, err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Println(c.errorMessage, err)
		return
	}

	countInBody := int64(bytes.Count(body, []byte(objectiveString)))

	atomic.AddInt64(totalValue, countInBody)

	res := fmt.Sprintf(c.outputPhrase, url, countInBody)

	outputChan <- res
}

// NewCounter creates new Counter interface
func NewCounter(httpClient *http.Client,
	errorMessage string,
	logger *log.Logger,
	outputPhrase string) *counter {
	return &counter{
		httpClient:   httpClient,
		errorMessage: errorMessage,
		logger:       logger,
		outputPhrase: outputPhrase,
	}
}
