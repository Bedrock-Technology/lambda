package core

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type FetchArgs struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type FetchResp struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func Fetch(url string, args FetchArgs) (*FetchResp, error) {
	slog.Debug("fetch()", "url", url, "args", args)

	var body io.Reader
	if args.Body != "" {
		body = io.NopCloser(strings.NewReader(args.Body))
	}

	req, err := http.NewRequest(args.Method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range args.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &FetchResp{
		Status: resp.StatusCode,
		Body:   string(respBody),
	}, nil
}
