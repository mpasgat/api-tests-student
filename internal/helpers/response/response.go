package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadBody проверяет статус-код ответа и возвращает тело.
// При несовпадении статуса тело попадает в сообщение об ошибке, чтобы было видно, что ответил сервер
func ReadBody(t *testing.T, resp *http.Response, expectedStatusCode int) []byte {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equalf(t, expectedStatusCode, resp.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	return body
}

// ReadJSON проверяет статус-код ответа и парсит тело в переданную структуру
func ReadJSON(t *testing.T, resp *http.Response, expectedStatusCode int, target any) {
	body := ReadBody(t, resp, expectedStatusCode)
	require.NoErrorf(t, json.Unmarshal(body, target), "Не удалось распарсить ответ: %s", string(body))
}
