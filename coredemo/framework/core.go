package framework

import "net/http"

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO
}
