package app

import sdk "github.com/davebryson/menta/types"

type route struct {
	handler sdk.TxHandler
}

type router struct {
	routes map[string]sdk.TxHandler
}

func NewRouter() *router {
	return &router{
		routes: make(map[string]sdk.TxHandler, 0),
	}
}

func (self *router) Add(path string, handler sdk.TxHandler) {
	self.routes[path] = handler
}

func (self *router) GetHandler(path string) sdk.TxHandler {
	return self.routes[path]
}
