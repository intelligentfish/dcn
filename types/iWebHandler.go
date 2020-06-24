package types

import "github.com/gin-gonic/gin"

// IWebHandler web handler
type IWebHandler interface {
	// UseEngine use web engine
	UseEngine(engine *gin.Engine)
}
