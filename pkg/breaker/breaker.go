package breaker

import (
	"context"
	"github.com/ahmadrezamusthafa/logwatcher/common/errors"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	cb "github.com/eapache/go-resiliency/breaker"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	CircuitBreakerOpenError = errors.New("Circuit breaker opened")
	NilHTTPClientError      = errors.New("HTTP client not initialized")
	HTTPResponseError       = errors.New("Internal server error")
	HTTPRequestError        = errors.New("Client request error")
)

type CbBreakerConfig struct {
	ErrorThreshold   int
	SuccessThreshold int
	Timeout          time.Duration
	HttpTimeout      time.Duration
}

type HttpClient struct {
	client  *http.Client
	breaker *cb.Breaker
}

type ClientBreaker struct {
	breaker *cb.Breaker
}

func initDefaultHTTPClient(httpTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

func NewHTTPClient(cbCfg CbBreakerConfig, client *http.Client) *HttpClient {
	cb := cb.New(cbCfg.ErrorThreshold, cbCfg.SuccessThreshold, cbCfg.Timeout)
	if client == nil {
		client = initDefaultHTTPClient(cbCfg.HttpTimeout)
	}

	httpClient := &HttpClient{
		client:  client,
		breaker: cb,
	}

	return httpClient
}

func (c *HttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c == nil {
		return nil, NilHTTPClientError
	}

	var resp *http.Response
	res := c.breaker.Run(func() error {
		var err error
		resp, err = c.client.Do(req)
		if resp != nil {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				jsonByte, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				if val := jsoniter.Get(jsonByte, "error", "message").ToString(); val != "" {
					return errors.New(val)
				}
				return HTTPRequestError
			} else if resp.StatusCode >= 500 {
				return HTTPResponseError
			}
		}
		return err
	})

	if res == cb.ErrBreakerOpen {
		return nil, CircuitBreakerOpenError
	}

	return resp, res
}

func NewClientBreaker(cbCfg CbBreakerConfig) *ClientBreaker {
	cb := cb.New(cbCfg.ErrorThreshold, cbCfg.SuccessThreshold, cbCfg.Timeout)
	breaker := &ClientBreaker{
		breaker: cb,
	}

	return breaker
}

func (c *ClientBreaker) Do(ctx context.Context, work func() error) error {
	cbErr := c.breaker.Run(work)
	if cbErr == cb.ErrBreakerOpen {
		return CircuitBreakerOpenError
	}

	return cbErr
}

type EngineBreaker struct {
	Config     config.Config `inject:"config"`
	HttpClient *HttpClient
	*ClientBreaker
}

func (mb *EngineBreaker) StartUp() {
	logger.Info("Init breaker... ")

	cbConfig := CbBreakerConfig{
		ErrorThreshold:   mb.Config.BreakerErrorThreshold,
		SuccessThreshold: mb.Config.BreakerSuccessThreshold,
		Timeout:          mb.Config.BreakerTimeout,
		HttpTimeout:      mb.Config.HttpClientTimeout,
	}

	mb.ClientBreaker = NewClientBreaker(cbConfig)
	mb.HttpClient = NewHTTPClient(cbConfig, nil)
}

func (mb *EngineBreaker) Shutdown() {}
