package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// mockDoer implements Doer for testing. It returns responses in order.
type mockDoer struct {
	responses []mockResponse
	callIdx   int
	calls     []mockCall
}

type mockResponse struct {
	status int
	body   any
}

type mockCall struct {
	method string
	path   string
	body   any
}

func (m *mockDoer) Do(_ context.Context, method, path string, body any) (*http.Response, error) {
	m.calls = append(m.calls, mockCall{method: method, path: path, body: body})

	if m.callIdx >= len(m.responses) {
		return nil, fmt.Errorf("no more mock responses (call %d)", m.callIdx)
	}
	r := m.responses[m.callIdx]
	m.callIdx++
	data, _ := json.Marshal(r.body)
	return &http.Response{
		StatusCode: r.status,
		Body:       io.NopCloser(bytes.NewReader(data)),
		Header:     http.Header{},
	}, nil
}
