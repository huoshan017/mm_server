package main

import (
	msg_client_message "mm_server_new/proto/gen_go/client_message"
	"sync"
)

// CatInfoPool ...
type CatInfoPool struct {
	pool *sync.Pool
}

// Init ...
func (p *CatInfoPool) Init() {
	p.pool = &sync.Pool{
		New: func() interface{} {
			return &msg_client_message.CatInfo{}
		},
	}
}

// Get ...
func (p *CatInfoPool) Get() *msg_client_message.CatInfo {
	return p.pool.Get().(*msg_client_message.CatInfo)
}

// Put ...
func (p *CatInfoPool) Put(ds *msg_client_message.CatInfo) {
	p.pool.Put(ds)
}
