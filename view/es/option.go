package es

import (
	"context"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/xxjwxc/public/mylog"
)

type options struct {
	retries int

	timeout time.Duration

	indexName string
	typeName  string
	addrs     []string
	ctx       context.Context
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithAddrs 设置地址
func WithAddrs(addrs ...string) Option {
	return optionFunc(func(o *options) {
		o.addrs = append(o.addrs, addrs...)
	})
}

// WithCtx ...
func WithCtx(ctx context.Context) Option {
	return optionFunc(func(o *options) {
		o.ctx = ctx
	})
}

// WithIndexName 设置索引
func WithIndexName(indexName string) Option {
	return optionFunc(func(o *options) {
		o.indexName = indexName
	})
}

// WithTypeName 设置类型
func WithTypeName(typeName string) Option {
	return optionFunc(func(o *options) {
		o.typeName = typeName
	})
}

// WithRetries 设置重试次数
func WithRetries(retries int) Option {
	return optionFunc(func(o *options) {
		o.retries = retries
	})
}

// WithTimeout 设置超时
func WithTimeout(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timeout = t
	})
}

//
func (es *MyElastic) newClient() error {
	var err error

	httpClient := &http.Client{
		Timeout: es.ops.timeout,
	}

	for i := 0; i < es.ops.retries; i++ {
		es.client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es.ops.addrs...), elastic.SetHttpClient(httpClient), elastic.SetMaxRetries(es.ops.retries))
		if err == nil {
			break
		}
		mylog.Error(err)
		time.Sleep(100 * time.Millisecond)
	}

	return err
}
