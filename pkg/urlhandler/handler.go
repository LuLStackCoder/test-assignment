package urlhandler

import (
	"bufio"
	"io"
	"sync"
)

// URLHandler interface represent the objective string counter from all pages read from reader
type URLHandler interface {
	CountAllUrls(reader io.Reader, objectiveString string) (count uint32)
}

type counter interface {
	Count(url string, wg *sync.WaitGroup, objectiveString string, value *uint32)
}

type urlHandler struct {
	counter counter
	wg      *sync.WaitGroup
}

// CountAllUrls count the objective string from all pages read from reader
func (h *urlHandler) CountAllUrls(reader io.Reader, objectiveString string) (count uint32) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		h.wg.Add(1)

		go h.counter.Count(scanner.Text(), h.wg, objectiveString, &count)
	}

	h.wg.Wait()

	return
}

// NewURLHandler creates new URLHandler interface
func NewURLHandler(counter counter, wg *sync.WaitGroup) URLHandler {
	return &urlHandler{
		counter: counter,
		wg:      wg,
	}
}
