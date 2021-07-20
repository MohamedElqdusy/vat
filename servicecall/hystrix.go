package servicecall

import (
	"fmt"
	"net/http"
	"time"
	"vat/logger"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
)

func retry(req *http.Request, ResponseChan chan *http.Response) error {
	client := http.Client{}
	r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
	attempt := 0
	err := r.Run(func() error {
		attempt++
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			if err == nil {
				ResponseChan <- resp
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("Status was %v", resp.StatusCode)
		}

		logger.Errorf("Retrier failed, attempt %v", attempt)

		return err
	})
	return err
}

// CircuitBreaker used for calling upstream services with circuitBreaker
func CircuitBreaker(breakerName string, req *http.Request) (*http.Response, error) {
	outputChan := make(chan *http.Response, 1)
	errors := hystrix.Go(breakerName, func() error {
		err := retry(req, outputChan)
		return err
	}, func(err error) error {
		logger.Errorf("Circute breaker %v has [Error]: %v", breakerName, err.Error())
		return err
	})

	select {
	case out := <-outputChan:
		logger.Infof("Breaker %v successful", breakerName)
		return out, nil

	case err := <-errors:
		logger.Errorf("Circute breaker %v has [Error]: %v", breakerName, err.Error())
		return nil, err
	}
}
