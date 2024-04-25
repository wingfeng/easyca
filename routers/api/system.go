package api

import (
	"easyca/controller"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

//后端api
func RegisterSystem(r *gin.RouterGroup) {
	store := persistence.NewInMemoryStore(time.Second)

	api := r.Group("/system")
	{
		api.Any("/current", controller.Current)

		api.Any("/version", controller.Version)

		api.GET("/menu", cache.CachePage(store, time.Hour, controller.Menu))
	}
}
