package server_worker

import (
	"cpm-standings/algocode_worker"
	"cpm-standings/config"
	"cpm-standings/data_worker"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func CraftStandings(conf *config.Config) *gin.H {
	criteria := config.ParseCriteria(conf.CriteriaPath)
	data := algocode_worker.GetSubmitsData(conf.SubmitsLink)
	res := data_worker.ExportStudentsData(data, criteria)
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

func RunServer(conf *config.Config) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.gohtml")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "standings.gohtml", CraftStandings(conf))
	})
	router.GET("/criteria", func(c *gin.Context) {
		c.HTML(http.StatusOK, "criteria.gohtml", CraftCriteria(conf))
	})

	err := router.Run(conf.Host)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
