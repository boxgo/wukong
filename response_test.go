package wukong_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boxgo/wukong"
)

func TestResponseIsTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	resp := wukong.New(ts.URL).Get("/").WithCTX(ctx).End()

	wukong.ExpectEqual(t, true, resp.IsTimeout())
}

func TestResponseIsCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 500)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cancel()
	}()

	resp := wukong.New(ts.URL).Get("/").WithCTX(ctx).End()

	wukong.ExpectEqual(t, true, resp.IsCancel())
}

func TestResponseBindJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type data struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}
	d := &data{}

	resp := wukong.New(ts.URL).Get("/").WithCTX(ctx).End().BindBody(d)

	wukong.ExpectEqual(t, resp.Error(), nil)
	wukong.ExpectEqual(t, d.String, "string")
	wukong.ExpectEqual(t, d.Int, 1)
	wukong.ExpectEqual(t, d.Float, 2.3)
	wukong.ExpectEqual(t, d.Bool, true)
}

func TestResponseConditionBindJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
		} else {
			w.Write([]byte(`{"errcode":1,"errmsg": "not ok"}`))
		}
	}))
	defer ts.Close()

	type Err struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	type Data struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	condition := func(d interface{}) (ok bool) {
		switch v := d.(type) {
		case *Err:
			ok = v.ErrCode != 0 && v.ErrMsg != ""
		case *Data:
			ok = true
		default:
			ok = false
		}

		return ok
	}

	{
		var (
			err  Err
			data Data
		)
		e := wukong.New(ts.URL).Get("/").End().ConditionBindBody(condition, &err, &data).Error()

		wukong.ExpectEqual(t, e, nil)
		wukong.ExpectEqual(t, err.ErrCode, 0)
		wukong.ExpectEqual(t, err.ErrMsg, "")
		wukong.ExpectEqual(t, data.Bool, true)
		wukong.ExpectEqual(t, data.String, "string")
		wukong.ExpectEqual(t, data.Int, 1)
		wukong.ExpectEqual(t, data.Float, 2.3)
	}

	{
		var (
			err  Err
			data Data
		)
		e := wukong.New(ts.URL).Post("/").End().ConditionBindBody(condition, &err, &data).Error()

		wukong.ExpectEqual(t, e, nil)
		wukong.ExpectEqual(t, err.ErrCode, 1)
		wukong.ExpectEqual(t, err.ErrMsg, "not ok")
		wukong.ExpectEqual(t, data.Bool, false)
		wukong.ExpectEqual(t, data.String, "")
		wukong.ExpectEqual(t, data.Int, 0)
		wukong.ExpectEqual(t, data.Float, 0.0)
	}
}
