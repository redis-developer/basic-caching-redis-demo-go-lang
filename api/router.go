package api

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/redis-developer/basic-caching-redis-demo-go-lang/controller"
	"log"
	"net/http"
	"time"
)

func router(publicPath string) http.Handler {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile(publicPath, true)))
	router.GET("/repos/:username", handlerRepos)
	return router

}

func response(c *gin.Context, data interface{}, err error) {
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func handlerRepos(c *gin.Context) {

	timeStart := time.Now()

	username := c.Param("username")
	repo, err := controller.Instance().GetRepo(username)
	if err != nil {
		response(c, `{"error":"could not read repo"}`, err)
		return
	} else {
		repo.Duration = time.Now().Sub(timeStart)
		if repo.Cached == false {
			err := controller.Instance().SetNonCachedDuration(username, repo.Duration)
			if err != nil {
				response(c, `{"error":"could not save duration"}`, err)
				return
			}
		} else {
			duration, err := controller.Instance().GetNonCachedDuration(username)
			if err != nil {
				response(c, `{"error":"could not get duration"}`, err)
				return
			}
			log.Println(duration, repo.Duration)
			repo.Faster = int(duration / repo.Duration)

		}

		c.Header("X-Powered-By", "golang gin")
		c.Header("X-Response-Time", fmt.Sprintf("%.3fms", float64(repo.Duration)/float64(time.Millisecond)))
		response(c, repo, nil)
	}

}
