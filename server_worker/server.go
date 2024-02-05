package server_worker

import (
	"cpm-standings/config"
	"cpm-standings/data_worker"
	codeforcesapi "github.com/MichailKon/codeforces-api"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func CraftStandings(conf *config.Config, mapping config.StudentsHandlesMapping) *gin.H {
	criteria := config.ParseCriteria(conf.CriteriaPath)
	session := codeforcesapi.NewCodeforcesSession(conf.ApiKey, conf.Secret)
	res := data_worker.ExportStudentsData(session, mapping, criteria)
	return &gin.H{
		"ContestTitles": res.ContestTitles,
		"Students":      res.Students,
	}
}

func CraftCriteria(conf *config.Config) *gin.H {
	criteria := config.ParseCriteria(conf.CriteriaPath)
	return &gin.H{
		"Criteria": criteria,
	}
}

func RunServer(conf *config.Config, mapping config.StudentsHandlesMapping) {
	memoryStore := persist.NewMemoryStore(10 * time.Minute)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.gohtml")
	router.GET("/", cache.CacheByRequestURI(memoryStore, 10*time.Minute), func(c *gin.Context) {
		c.HTML(http.StatusOK, "standings.gohtml", CraftStandings(conf, mapping))
	})
	router.GET("/criteria", cache.CacheByRequestURI(memoryStore, 10*time.Minute), func(c *gin.Context) {
		c.HTML(http.StatusOK, "criteria.gohtml", CraftCriteria(conf))
	})

	err := router.Run(conf.Host)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
