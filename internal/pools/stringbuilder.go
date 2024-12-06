package pools

import (
	"strings"
	"sync"
)

type StringBuilder struct {
	pool sync.Pool
}

func NewStringBuilderPool() *StringBuilder {
	return &StringBuilder{}
}

func (sbp *StringBuilder) Get() *strings.Builder {
	v := sbp.pool.Get()
	if v == nil {
		return &strings.Builder{}
	}
	return v.(*strings.Builder)
}

func (sbp *StringBuilder) Put(ua *strings.Builder) {
	ua.Reset()
	sbp.pool.Put(ua)
}
