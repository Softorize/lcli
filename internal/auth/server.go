// server.go runs a temporary local HTTP server to receive the OAuth callback.
package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const successHTML = `<!DOCTYPE html>
<html><head><title>lcli</title></head>
<body style="font-family:sans-serif;text-align:center;padding:3rem">
<h2>Authentication successful!</h2>
<p>You can close this tab.</p>
</body></html>`

// callbackResult holds the values extracted from the OAuth redirect.
type callbackResult struct {
	Code  string
	State string
	Err   error
}

// CallbackServer is a temporary HTTP server that listens for the OAuth
// callback from LinkedIn.
type CallbackServer struct {
	port   int
	result chan callbackResult
}

// NewCallbackServer creates a CallbackServer that will listen on the given port.
func NewCallbackServer(port int) *CallbackServer {
	return &CallbackServer{
		port:   port,
		result: make(chan callbackResult, 1),
	}
}

// Start begins listening for the OAuth callback. It blocks until a callback is
// received or the context is cancelled. On success it returns the authorization
// code and state parameter.
func (s *CallbackServer) Start(ctx context.Context) (code string, state string, err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return "", "", fmt.Errorf("listen on port %d: %w", s.port, err)
	}

	go func() { _ = srv.Serve(ln) }()

	defer func() { _ = srv.Shutdown(context.Background()) }()

	select {
	case <-ctx.Done():
		return "", "", fmt.Errorf("callback server: %w", ctx.Err())
	case res := <-s.result:
		return res.Code, res.State, res.Err
	}
}

// handleCallback processes the OAuth redirect request.
func (s *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if errMsg := q.Get("error"); errMsg != "" {
		desc := q.Get("error_description")
		s.result <- callbackResult{Err: fmt.Errorf("oauth error: %s — %s", errMsg, desc)}
		http.Error(w, "Authentication failed.", http.StatusBadRequest)
		return
	}

	code := q.Get("code")
	if code == "" {
		s.result <- callbackResult{Err: fmt.Errorf("missing code parameter")}
		http.Error(w, "Missing code.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, successHTML)

	s.result <- callbackResult{Code: code, State: q.Get("state")}
}
