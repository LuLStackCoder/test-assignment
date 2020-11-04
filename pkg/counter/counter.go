package counter

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

// Counter interface represent the objective string counter on the single page
type Counter interface {
	Count(url string, wg *sync.WaitGroup, objectiveString string, value *uint32)
}

type counter struct {
	httpClient   *http.Client
	quotaChannel chan struct{}
	errorMessage string
	logger       *log.Logger
	output       io.Writer
	outputPhrase string
}

// Count objective string on the single page
func (c *counter) Count(url string, wg *sync.WaitGroup, objectiveString string, value *uint32) {
	c.quotaChannel <- struct{}{}

	defer func() {
		wg.Done()
		<-c.quotaChannel
	}()

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

	countInBody := uint32(bytes.Count(body, []byte(objectiveString)))

	atomic.AddUint32(value, countInBody)

	_, err = fmt.Fprintf(c.output, c.outputPhrase, url, countInBody)
	if err != nil {
		c.logger.Println(c.errorMessage, err)
		return
	}
}

// NewCounter creates new Counter interface
func NewCounter(httpClient *http.Client,
	quotaChannel chan struct{},
	errorMessage string,
	logger *log.Logger,
	output io.Writer,
	outputPhrase string) Counter {
	return &counter{
		httpClient:   httpClient,
		quotaChannel: quotaChannel,
		errorMessage: errorMessage,
		logger:       logger,
		output:       output,
		outputPhrase: outputPhrase,
	}
}
