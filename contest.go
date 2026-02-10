package contest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type T interface {
	Helper()
	Fatalf(format string, args ...any)
	Logf(format string, args ...any)
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

	ce, _ := detectConnErr(c.resp.Body.Bytes())
	if ce != nil {
		c.lastErr = ce
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

func (c *Client) Err() *connect.Error {
	c.t.Helper()
	c.ensureDid()
	if c.lastErr == nil {
		return nil
	}
	return c.lastErr
}

// ----- Helpers -----

func (c *Client) ensureDid() {
	c.t.Helper()
	if !c.didDo {
		c.t.Fatalf("call Do() before assertions")
	}
}

func detectConnErr(body []byte) (*connect.Error, error) {
	var rc rawConnErr
	if err := json.Unmarshal(body, &rc); err != nil {
		return nil, fmt.Errorf("unmarshal connect error: %w", err)
	}
	ce, err := rc.toConnErr()
	if err != nil {
		return nil, fmt.Errorf("convert to connect error: %w", err)
	}
	if len(rc.Details) == 0 {
		return ce, nil
	}
	details, err := rc.unmarshalDetails()
	if err != nil {
		return nil, fmt.Errorf("unmarshal error details: %w", err)
	}
	for _, d := range details {
		ce.AddDetail(d)
	}
	return ce, nil
}

type rawConnErr struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details []rawDetail `json:"details"`
}

func (rc *rawConnErr) toConnErr() (*connect.Error, error) {
	code, err := sToConnectCode(rc.Code)
	if err != nil {
		return nil, fmt.Errorf("invalid connect error code: %w", err)
	}
	return connect.NewError(code, errors.New(rc.Message)), nil
}

func (rc *rawConnErr) unmarshalDetails() ([]*connect.ErrorDetail, error) {
	details := make([]*connect.ErrorDetail, 0, len(rc.Details))
	for _, d := range rc.Details {
		if d.Value == "" {
			continue
		}
		raw, err := decodeB64(d.Value)
		if err != nil {
			return nil, fmt.Errorf("base64 decode detail: %w", err)
		}
		typeURL := d.Type
		if !strings.Contains(d.Type, "/") {
			typeURL = "type.googleapis.com/" + d.Type
		}
		ed, err := connect.NewErrorDetail(&anypb.Any{
			TypeUrl: typeURL,
			Value:   raw,
		})
		if err != nil {
			return nil, fmt.Errorf("create error detail: %w", err)
		}
		details = append(details, ed)
	}
	return details, nil
}

func decodeB64(s string) ([]byte, error) {
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return base64.RawStdEncoding.DecodeString(s)
}

type rawDetail struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
