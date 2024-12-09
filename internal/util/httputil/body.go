package httputil

import (
	"io"
	"net/http"
	"strings"

	"github.com/zNoah-1/Arc-MS/internal/logger"
)

func BodyBytes(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	var bodyBytes []byte
	var err error
	bodyBytes, err = io.ReadAll(r.Body)

	if err == nil {
		return bodyBytes, nil //All is super :)
	}

	logger.Error("Something wrong happened: ", err) //All is not okay ):
	return nil, err
}

// Returns empty string or "" (no quotes) if there's no value
func ResponseValue(bodyBytes []byte, key string) string {
	if len(bodyBytes) > 0 {
		bodyData := strings.Split(string(bodyBytes), "&")

		i := 0
		for i < len(bodyData) {
			dataEntry := strings.Split(bodyData[i], "=")
			entryKey := dataEntry[0]

			if entryKey == key {
				value := dataEntry[1]
				return value
			}
			i++
		}
	}
	return ""
}
