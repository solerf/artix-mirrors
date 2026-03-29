package mirror

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Server struct {
	Country string
	Url     string
	measure time.Duration
}

func (s *Server) measureConnection(mirrorTimeout time.Duration) error {
	ctx, cancel := context.WithTimeoutCause(
		context.Background(),
		mirrorTimeout,
		fmt.Errorf("connection timeout [%v] exceeded", mirrorTimeout),
	)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.Url, http.NoBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	var reqStart, reqEnd time.Time
	var respStart, respEnd time.Time

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			// the start
			reqStart = time.Now()
		},
		GotConn: func(gotConnInfo httptrace.GotConnInfo) {
			// the end of conn, I went through all steps... dns, conn start, tls handshake etc
			reqEnd = time.Now()
		},
		GotFirstResponseByte: func() {
			// response start
			respStart = time.Now()
		},
	}

	tracedRequest := request.WithContext(httptrace.WithClientTrace(request.Context(), trace))

	response, err := http.DefaultClient.Do(tracedRequest)
	if err != nil {
		return fmt.Errorf("call request: %w", err)
	}
	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	respEnd = time.Now()

	s.measure = reqEnd.Sub(reqStart) + respEnd.Sub(respStart)
	return nil
}
