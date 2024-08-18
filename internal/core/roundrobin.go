package core

import "sync/atomic"

type RoundRobin[T any] struct {
	index   int32
	objsLen int32
	objs    []*T
}

func NewRoundRobin[T any](objs ...*T) *RoundRobin[T] {
	return &RoundRobin[T]{
		index:   0,
		objs:    objs,
		objsLen: int32(len(objs)),
	}
}

func (r *RoundRobin[T]) Next() *T {
	if r.objsLen == 0 {
		return nil
	}
	if r.objsLen == 1 {
		return r.objs[0]
	}
	index := r.index
	for {
		index = atomic.LoadInt32(&r.index)
		newVal := (index + 1) % r.objsLen
		if atomic.CompareAndSwapInt32(&r.index, index, newVal) {
			break
		}
	}
	return r.objs[index]
}
