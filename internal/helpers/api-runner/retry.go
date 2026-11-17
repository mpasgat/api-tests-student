package api_runner

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	maxAttempts = 3
	retryDelay  = time.Second
)

// retryTransport повторяет запрос, если тестовый стенд ответил 503.
// Стенд периодически отдаёт 503 под нагрузкой, это не связано с проверяемой логикой
type retryTransport struct {
	next http.RoundTripper
}

// RoundTrip выполняет запрос и повторяет его до maxAttempts раз, пока сервер отвечает 503
func (rt retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body.Close()
	}

	for attempt := 1; ; attempt++ {
		if body != nil {
			request.Body = io.NopCloser(bytes.NewReader(body))
		}

		response, err := rt.next.RoundTrip(request)
		if err != nil || response.StatusCode != http.StatusServiceUnavailable || attempt == maxAttempts {
			return response, err
		}

		response.Body.Close()
		log.Printf("Сервер ответил 503 на %s %s, повтор через %s", request.Method, request.URL.Path, retryDelay)
		time.Sleep(retryDelay)
	}
}
