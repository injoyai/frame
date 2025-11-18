package gins

import "github.com/gin-gonic/gin"

type Grouper struct {
	*gin.RouterGroup
}

func (this *Grouper) Group(path string, f func(g *Grouper)) {
	group := this.RouterGroup.Group(path)
	f(&Grouper{RouterGroup: group})
}

func (this *Grouper) Use(middle ...Handler) {
	this.RouterGroup.Use(middle...)
}

func (this *Grouper) ALL(path string, handler Handler) {
	this.RouterGroup.Any(path, handler)
}
