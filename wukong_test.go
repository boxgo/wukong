package wukong_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boxgo/wukong"
)

func TestTrace(t *testing.T) {
	client := wukong.New("https://www.baidu.com")

	var (
		err        error
		statusCode int
		statusMsg  string
		header     http.Header
		isTimeout  bool
		isCancel   bool
		data       []byte
	)
	req := client.Get("/").WithCTX(context.Background())
	req.End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindError(&err).
		BindBytes(&data)

	t.Logf("%s \t DNSStart", req.TraceInfo.DNSStart)
	t.Logf("%s \t DNSDone", req.TraceInfo.DNSDone)
	t.Logf("%s \t ConnectStart", req.TraceInfo.ConnectStart)
	t.Logf("%s \t ConnectDone", req.TraceInfo.ConnectDone)
	t.Logf("%s \t GetConn", req.TraceInfo.GetConn)
	t.Logf("%s \t GotConn", req.TraceInfo.GotConn)
	t.Logf("%s \t TLSHandshakeStart", req.TraceInfo.TLSHandshakeStart)
	t.Logf("%s \t TLSHandshakeDone", req.TraceInfo.TLSHandshakeDone)
	t.Logf("%s \t GotFirstResponseByte", req.TraceInfo.GotFirstResponseByte)
	t.Log("GotConnReused ", req.TraceInfo.GotConnReused)
	t.Log("GotConnWasIdle ", req.TraceInfo.GotConnWasIdle)
	t.Log("GotConnIdleTime ", req.TraceInfo.GotConnIdleTime)

}

func TestSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	client := wukong.New("")

	type Body struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	var (
		statusCode int
		statusMsg  string
		bodyData   Body
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	err := client.Get(ts.URL).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindBody(&bodyData).
		Error()

	wukong.ExpectEqual(t, err, nil)
	wukong.ExpectEqual(t, statusCode, 200)
	wukong.ExpectEqual(t, statusMsg, "200 OK")
	wukong.ExpectEqual(t, isTimeout, false)
	wukong.ExpectEqual(t, isCancel, false)
}

func TestSample(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	client := wukong.New(ts.URL)

	type Body struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	var (
		statusCode int
		statusMsg  string
		bodyData   Body
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	err := client.Get("/").WithCTX(context.Background()).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindBody(&bodyData).
		Error()

	wukong.ExpectEqual(t, err, nil)
	wukong.ExpectEqual(t, statusCode, 200)
	wukong.ExpectEqual(t, statusMsg, "200 OK")
	wukong.ExpectEqual(t, isTimeout, false)
	wukong.ExpectEqual(t, isCancel, false)
}
