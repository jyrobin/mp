package mpi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jyrobin/goutil"
	"github.com/jyrobin/mp"
)

type Request struct {
	Method  string        `json:"method,omitempty"`
	Meta    mp.MetaJson   `json:"meta,omitempty"`
	Options []mp.MetaJson `json:"options,omitempty"` //TODO: rename to Options
}

func (req Request) Unpack() (string, mp.Meta, []mp.Meta) {
	m, opts := mp.JsonToMeta(req.Meta), mp.JsonsToMetas(req.Options)
	return req.Method, m, opts
}

func RequestBody(method string, m mp.Meta, opts []mp.Meta) []byte {
	req := Request{method, mp.MetaToJson(m), mp.MetasToJsons(opts)}
	buf, _ := json.MarshalIndent(req, "", "  ")
	return buf
}

type LocalMpi struct {
	srv    http.Handler
	prefix string
}

func NewLocalMpi(srv http.Handler, prefix string) LocalMpi {
	return LocalMpi{srv, prefix}
}

func NewRemoteMpi(hostPort, prefix string) (*LocalMpi, error) {
	if srv, err := goutil.UriHandler(hostPort, nil); err != nil {
		return nil, err
	} else {
		return &LocalMpi{srv, prefix}, nil
	}
}

func (mpi LocalMpi) IsNil() bool {
	return goutil.IsNil(mpi.srv)
}

func (mpi LocalMpi) Prefix() string {
	return mpi.prefix
}

func (mpi LocalMpi) Handler() http.Handler {
	return mpi.srv
}

func (mpi LocalMpi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mpi.srv.ServeHTTP(w, req)
}

func (mpi LocalMpi) Call(ctx context.Context, method string, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	uri := fmt.Sprintf("%s/call", mpi.prefix)
	body := RequestBody(method, m, opts)
	req, err := goutil.AjaxRequest("POST", uri, body, nil)
	if err != nil {
		return mp.Nil, err
	}

	res := goutil.NewResponseWriter()
	mpi.srv.ServeHTTP(res, req)

	var mj mp.MetaJson
	if err := res.Unmarshal(&mj, true); err != nil {
		return mp.Nil, err
	}
	return mp.JsonToMeta(mj), nil
}

func (mpi LocalMpi) List(ctx context.Context, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	return mpi.Call(ctx, "list", m, opts...)
}

func (mpi LocalMpi) Find(ctx context.Context, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	return mpi.Call(ctx, "find", m, opts...)
}

func (mpi LocalMpi) Create(ctx context.Context, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	return mpi.Call(ctx, "create", m, opts...)
}

func (mpi LocalMpi) Make(ctx context.Context, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	return mpi.Call(ctx, "make", m, opts...)
}

/* Later
func (mpi localMpi) List(ctx context.Context, m mp.Meta, opts ...mp.Meta) ([]mp.Meta, error) {
	uri := fmt.Sprintf("%s/list", mpi.prefix)
	body := RequestBody("list", m, opts)
	req, err := goutil.AjaxRequest("POST", uri, body, nil)
	if err != nil {
		return nil, err
	}

	res := goutil.NewResponseWriter()
	mpi.handler.ServeHTTP(res, req)

	var mjs []mp.MetaJson
	if err := res.Unmarshal(&mjs); err != nil {
		return nil, err
	}

	return mp.JsonsToMetas(mjs), nil
}

func (mpi localMpi) First(ctx context.Context, m mp.Meta, opts ...mp.Meta) (mp.Meta, error) {
	uri := fmt.Sprintf("%s/first", mpi.prefix)
	body := RequestBody("first", m, opts)
	req, err := goutil.AjaxRequest("POST", uri, body, nil)
	if err != nil {
		return nil, err
	}

	res := goutil.NewResponseWriter()
	mpi.handler.ServeHTTP(res, req)

	var mj mp.MetaJson
	if err := res.Unmarshal(&mj); err != nil {
		return nil, err
	}

	return mp.JsonToMeta(mj), nil
}
*/
