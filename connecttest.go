package connecttest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

type T interface {
	Helper()
	Fatalf(format string, args ...any)
}

type Client struct {
	t           T
	h           http.Handler
	proc        string
	headers     http.Header
	body        []byte
	contentType string
	resp        *httptest.ResponseRecorder
	lastErr     *connect.Error
	didDo       bool
}

func New(t T, h http.Handler) *Client {
	t.Helper()
	return &Client{
		t:       t,
		h:       h,
		headers: make(http.Header),
	}
}

// ----- Request configuration -----

func (c *Client) Procedure(proc string) *Client {
	c.proc = proc
	return c
}

func (c *Client) Header(key, value string) *Client {
	c.headers.Add(key, value)
	return c
}

func (c *Client) In(msg proto.Message) *Client {
	c.t.Helper()
	b, err := proto.Marshal(msg)
	if err != nil {
		c.t.Fatalf("marshal proto: %v", err)
	}
	c.body = b
	c.contentType = "application/proto"
	return c
}

// ----- Execute request -----

func (c *Client) Do() *Client {
	c.t.Helper()
	if c.proc == "" {
		c.t.Fatalf("procedure not set")
	}
	req := httptest.NewRequest(http.MethodPost, c.proc, bytes.NewReader(c.body))
	if c.contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}
	for k, vs := range c.headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	recorder := httptest.NewRecorder()
	c.h.ServeHTTP(recorder, req)
	c.resp = recorder
	c.didDo = true

	if code, ok := connectCodeFromHeaders(recorder.Result()); ok {
		c.lastErr = connect.NewError(code, errors.New(recorder.Body.String()))
	}

	return c
}

// ----- Response assertions -----

func (c *Client) ExpectStatus(code int) *Client {
	c.t.Helper()
	c.ensureDid()
	if c.resp.Code != code {
		c.t.Fatalf("status code: got %d want %d", c.resp.Code, code)
	}
	return c
}

func (c *Client) ExpectHeader(key string, want ...string) *Client {
	c.t.Helper()
	c.ensureDid()
	got := c.resp.Header().Values(key)
	if len(want) == 0 {
		if len(got) == 0 {
			c.t.Fatalf("header %q not present", key)
		}
		return c
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		c.t.Fatalf("header %q: got %v want %v", key, got, want)
	}
	return c
}

func (c *Client) Out(dst proto.Message) *Client {
	c.t.Helper()
	err := proto.Unmarshal(c.resp.Body.Bytes(), dst)
	if err != nil {
		c.t.Fatalf("unmarshal proto: %v; body=%s", err, c.resp.Body.String())
	}
	return c
}

// ----- Helpers -----

func (c *Client) ensureDid() {
	c.t.Helper()
	if !c.didDo {
		c.t.Fatalf("call Do() before assertions")
	}
}

func connectCodeFromHeaders(r *http.Response) (connect.Code, bool) {
	if v := r.Header.Get("Grpc-Status"); v != "" {
		connectCode(v)
	}
	return connect.CodeUnknown, false
}

func connectCode(s string) (connect.Code, bool) {
	code, err := strconv.Atoi(s)
	if err != nil {
		return connect.CodeUnknown, false
	}
	return connect.Code(code), true
}
