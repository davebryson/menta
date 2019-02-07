package app

import sdk "github.com/davebryson/menta/types"

type route struct {
	path    string
	handler sdk.TxHandler
}

type router struct {
	routes []route
}

func NewRouter() *router {
	return &router{
		routes: make([]route, 0),
	}
}

func (self *router) Add(path string, handler sdk.TxHandler) {
	self.routes =
		append(self.routes, route{path: path, handler: handler})
}

func (self *router) GetHandler(path string) sdk.TxHandler {
	for _, h := range self.routes {
		if h.path == path {
			return h.handler
		}
	}
	return nil
}
