package wukong

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

type (
	TraceInfo struct {
		DNSStart             time.Time
		DNSDone              time.Time
		ConnectStart         time.Time
		ConnectDone          time.Time
		GetConn              time.Time
		GotConn              time.Time
		GotConnReused        bool
		GotConnWasIdle       bool
		GotConnIdleTime      time.Duration
		TLSHandshakeStart    time.Time
		TLSHandshakeDone     time.Time
		GotFirstResponseByte time.Time
	}
)

func traceGenerator(request *Request) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			request.TraceInfo.DNSStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			request.TraceInfo.DNSDone = time.Now()
		},
		ConnectStart: func(_, _ string) {
			request.TraceInfo.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			request.TraceInfo.ConnectDone = time.Now()
		},
		GetConn: func(_ string) {
			request.TraceInfo.GetConn = time.Now()
		},
		GotConn: func(info httptrace.GotConnInfo) {
			request.TraceInfo.GotConn = time.Now()
			request.TraceInfo.GotConnReused = info.Reused
			request.TraceInfo.GotConnWasIdle = info.WasIdle
			request.TraceInfo.GotConnIdleTime = info.IdleTime
		},
		TLSHandshakeStart: func() {
			request.TraceInfo.TLSHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			request.TraceInfo.TLSHandshakeDone = time.Now()
		},
		GotFirstResponseByte: func() {
			request.TraceInfo.GotFirstResponseByte = time.Now()
		},
	}
}
