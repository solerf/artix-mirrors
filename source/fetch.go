package source

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetchMirrorList(url string) (io.Reader, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFunc()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return nil, fmt.Errorf("bad response: %s", response.Status)
	}

	buff := bytes.NewBuffer(nil)
	_, err = io.Copy(buff, response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return buff, nil
}
