package commonhandlers

import (
	"context"
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/common/respwriter"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	newrelic "github.com/newrelic/go-agent"

	uuid "github.com/satori/go.uuid"
)

// CommonHandlers is a collection of usefull middleware that could be used with http default package and gorilla mux
type CommonHandlers struct {
	nrApp newrelic.Application
	r     *mux.Router
}

// New will create an empty CommonHandlers
func New() *CommonHandlers {
	return &CommonHandlers{
		r: mux.NewRouter(),
	}
}

// SetNewRelic will setup new relic middleware utility and must be called if you want to use MonitoringHandler
func (c *CommonHandlers) SetNewRelic(applicationName string, key string, enabled bool) {
	nrConf := newrelic.NewConfig(applicationName, key)
	nrConf.Enabled = enabled
	nrApp, err := newrelic.NewApplication(nrConf)
	c.nrApp = nrApp
	if err != nil {
		logger.Err("New relic error: %v", err)
	}
}

func convertToNamespace(applicationName string) string {
	replacer := strings.NewReplacer(" ", "_",
		"(", "",
		")", "",
		"*", "",
		"!", "",
		"@", "",
		"#", "",
		"$", "",
		"%", "")
	namespace := replacer.Replace(applicationName)
	namespace = strings.ToLower(namespace)
	return namespace
}

// HeaderUserIDValidationHandler is helpful for making sure there is an Salestock UserID passed in header request
func (c *CommonHandlers) HeaderUserIDValidationHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		frags := strings.Split(r.URL.Path, "/")
		lastFrag := frags[len(frags)-1]
		if lastFrag != "health" {
			userID, err := uuid.FromString(r.Header.Get("UserID"))
			if err != nil {
				ssUserID, err := uuid.FromString(r.Header.Get("X-SS-User-ID"))
				if err != nil {
					logger.Warn("Invalid header: %v/%v", userID, ssUserID)
					respWriter := respwriter.New()
					respWriter.ErrorWriter(w, http.StatusUnauthorized, "en", nil)
					return
				}
				ctx := context.WithValue(r.Context(), "X-SS-User-ID", ssUserID)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Deprecated: use the new commonhandlers.Handle instead
// MonitoringHandler will send transaction time to NewRelic
func (c *CommonHandlers) MonitoringHandler(next http.Handler) http.Handler {
	if c.nrApp != nil {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, _ := url.ParseRequestURI(r.RequestURI)

			if u.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			txn := c.nrApp.StartTransaction(r.Method+":"+u.Path, w, r)
			defer txn.End()
			next.ServeHTTP(txn, r)
		}
		return http.HandlerFunc(fn)
	}
	return next
}

// Routes will return router instance
func (c *CommonHandlers) Routes() *mux.Router {
	return c.r
}

// Handle will handle current path with specified handler
func (c *CommonHandlers) Handle(method string, path string, handler func(w http.ResponseWriter, r *http.Request)) {
	// commented untul build is fixed
	// c.r.HandleFunc(path, prometheus.InstrumentHandlerFunc(path, c.newRelicWrapper(handler))).Methods(method)
	c.r.HandleFunc(path, c.newRelicWrapper(handler)).Methods(method)

}

func (c *CommonHandlers) newRelicWrapper(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	if c.nrApp == nil {
		return next
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		p, _ := mux.CurrentRoute(r).GetPathTemplate()

		if p == "/health" {
			next(w, r)
			return
		}

		txn := c.nrApp.StartTransaction(r.Method+" "+p, w, r)
		defer txn.End()
		next(txn, r)
	}
	return fn
}

// responseWriterDelegator to delegate the current writer
// this is a 100% from prometheus delegator with some modification
// the modification is needed because namespace is required
type responseWriterDelegator struct {
	http.ResponseWriter
	// handler, method string
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeStatusCode(status int) string {
	code := strconv.Itoa(status)
	return code
}

// RecoverHandler will catch and try to recover panic when processing any request
func (c *CommonHandlers) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Err("Panic: %v", err)
				respWriter := respwriter.New()
				respWriter.ErrorWriter(w, http.StatusInternalServerError, "en", fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// LoggingHandler will log every request and response to papertrail including method, path, response code, and time elapsed
func (c *CommonHandlers) LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(r.Header.Get("x-request-id"))
		if err != nil {
			id = uuid.NewV4()
		}

		ctx := context.WithValue(r.Context(), "x-request-id", id.String())
		r = r.WithContext(ctx)

		logger.Info("Request #%s - [%s] %q", id.String(), r.Method, r.URL.Path)
		start := time.Now()

		resp := httptest.NewRecorder()
		next.ServeHTTP(resp, r)

		for k, v := range resp.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.Code)
		s := resp.Body.String()
		code := resp.Code / 1e2
		elapsed := time.Since(start)
		switch code {
		case 2:
			logger.Info("Response #%s - [%s] %q - Result with status - %v - took %s", id.String(), r.Method, r.URL.Path, resp.Code, elapsed)
		case 4:
			logger.Warn("Response #%s - [%s] %q -  %s - Failed with status - %v - took %s", id.String(), r.Method, r.URL.Path, s, resp.Code, elapsed)
		case 5:
			logger.Warn("Response #%s - [%s] %q - %s - Failed with status - %v - took %s", id.String(), r.Method, r.URL.Path, s, resp.Code, elapsed)
		default:
			logger.Info("Response #%s - [%s] %q - %s - with status - %v - took %s", id.String(), r.Method, r.URL.Path, s, resp.Code, elapsed)
		}
		resp.Body.WriteTo(w)
	}

	return http.HandlerFunc(fn)
}

// Chain would create chained http handlers
func Chain(handlers ...alice.Constructor) alice.Chain {
	return alice.New(handlers...)
}
