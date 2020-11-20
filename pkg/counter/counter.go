package counter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)

// Сounter implements Counter
type Сounter struct {
	httpClient   *http.Client
	errorMessage string
	logger       *log.Logger
	outputPhrase string
}

// Count objective string on the single page
func (c *Сounter) Count(url string, objectiveString string, totalValue *int64, outputChan chan string) {
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

// NewCounter constructs and returns a usable Counter object
func NewCounter(httpClient *http.Client,
	errorMessage string,
	logger *log.Logger,
	outputPhrase string) *Сounter {
	return &Сounter{
		httpClient:   httpClient,
		errorMessage: errorMessage,
		logger:       logger,
		outputPhrase: outputPhrase,
	}
}
