package wukong

import (
	"net/http"
	"time"
)

type (
	WuKong struct {
		baseUrl   string
		client    *http.Client
		basicAuth BasicAuth
		query     map[string]string
		header    map[string]string
	}

	BasicAuth struct {
		Username string
		Password string
	}
)

func New(baseUrl string) *WuKong {
	w := &WuKong{
		baseUrl: baseUrl,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 60,
		},
	}

	return w
}

func (wk *WuKong) SetTimeout(t time.Duration) *WuKong {
	wk.client.Timeout = t

	return wk
}

func (wk *WuKong) SetTransport(transport *http.Transport) *WuKong {
	wk.client.Transport = transport

	return wk
}

func (wk *WuKong) SetBasicAuth(auth BasicAuth) *WuKong {
	wk.basicAuth = auth

	return wk
}

func (wk *WuKong) SetQuery(query map[string]string) *WuKong {
	wk.query = query

	return wk
}

func (wk *WuKong) SetHeader(header map[string]string) *WuKong {
	wk.header = header

	return wk
}

func (wk *WuKong) Get(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodGet, path))
}

func (wk *WuKong) Post(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPost, path))
}

func (wk *WuKong) Put(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPut, path))
}

func (wk *WuKong) Patch(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPatch, path))
}

func (wk *WuKong) Delete(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodDelete, path))
}

func (wk *WuKong) Head(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodHead, path))
}

func (wk *WuKong) Options(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodOptions, path))
}

func (wk *WuKong) Trace(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodTrace, path))
}

func (wk *WuKong) Client() *http.Client {
	return wk.client
}

func (wk *WuKong) initRequest(request *Request) *Request {
	request.Query(wk.query)
	request.SetBasicAuth(wk.basicAuth)

	for k, v := range wk.header {
		request.SetHeader(k, v)
	}

	return request
}

func (wk *WuKong) do(req *Request) (resp *Response) {
	var (
		err     error
		rawReq  *http.Request
		rawResp *http.Response
	)

	for ok := true; ok; ok = !ok {
		rawReq, err = req.RawRequest()
		if err != nil {
			break
		}

		if rawResp, err = wk.client.Do(rawReq); err != nil {
			break
		}
	}

	resp = NewResponse(err, req, rawResp)

	return resp
}
