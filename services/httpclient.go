package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type CustomTransport struct {
	Base         http.RoundTripper
	trimPrefix   string
	DebugLogPath string
	logMu        sync.Mutex
}

func NewCustomTransport(proxyURL string) *CustomTransport {
	transport := &http.Transport{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse PROXY_URL: %v", err))
		}
		transport.Proxy = http.ProxyURL(proxy)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	debugLogPath := os.Getenv("ATLASSIAN_DEBUG")
	trimPrefix := os.Getenv("ATLASSIAN_TRIM")
	return &CustomTransport{
		Base:         transport,
		trimPrefix:   trimPrefix,
		DebugLogPath: debugLogPath,
	}
}

func (t *CustomTransport) writeDebugLog(msg string) {
	t.logMu.Lock()
	defer t.logMu.Unlock()
	f, err := os.OpenFile(t.DebugLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(msg)
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	if len(t.trimPrefix) > 0 && strings.HasPrefix(req.URL.Path, t.trimPrefix) {
		req.URL.Path = strings.Replace(req.URL.Path, t.trimPrefix, "", 1)
	}
	if len(t.DebugLogPath) > 0 {
		reqBody, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(io.NopCloser(io.MultiReader(bytes.NewReader(reqBody))))
		msg := "\n--- HTTP REQUEST ---\n"
		msg += time.Now().Format(time.RFC3339Nano) + "\n"
		msg += req.Method + " " + req.URL.String() + "\n"
		for k, v := range req.Header {
			msg += k + ": " + fmt.Sprint(v) + "\n"
		}
		msg += "Body: " + string(reqBody) + "\n"
		t.writeDebugLog(msg)
	}

	resp, err := t.Base.RoundTrip(req)

	if len(t.DebugLogPath) > 0 && resp != nil {
		respBody, _ := io.ReadAll(resp.Body)
		msg := "\n--- HTTP RESPONSE ---\n"
		msg += time.Now().Format(time.RFC3339Nano) + "\n"
		msg += "Status: " + resp.Status + "\n"
		for k, v := range resp.Header {
			msg += k + ": " + fmt.Sprint(v) + "\n"
		}
		msg += "Body: " + string(respBody) + "\n"
		msg += "Duration: " + time.Since(start).String() + "\n"
		t.writeDebugLog(msg)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}
	return resp, err
}

var DefaultHttpClient = sync.OnceValue(func() *http.Client {
	proxyURL := os.Getenv("PROXY_URL")
	return &http.Client{Transport: NewCustomTransport(proxyURL)}
})
