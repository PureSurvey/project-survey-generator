package pools

import (
	"strings"
	"sync"
)

type StringBuilderPool struct {
	pool sync.Pool
}

func NewStringBuilderPool() *StringBuilderPool {
	return &StringBuilderPool{}
}

func (sbp *StringBuilderPool) Get() *strings.Builder {
	v := sbp.pool.Get()
	if v == nil {
		return &strings.Builder{}
	}
	return v.(*strings.Builder)
}

func (sbp *StringBuilderPool) Put(ua *strings.Builder) {
	ua.Reset()
	sbp.pool.Put(ua)
}
