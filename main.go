package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    initDB()
    defer db.Close()

    r := gin.Default()
    r.Use(loggingMiddleware())

    r.POST("/cats", createCat)
    r.GET("/cats", getCats)
    r.GET("/cats/:id", getCat)
    r.PUT("/cats/:id", updateCat)
    r.DELETE("/cats/:id", deleteCat)

    r.POST("/missions", createMission)
    r.GET("/missions", getMissions)
    r.GET("/missions/:id", getMission)
    r.PUT("/missions/:id", updateMission)
    r.DELETE("/missions/:id", deleteMission)

    r.POST("/targets/:mission_id", addTarget)
    r.PUT("/targets/:id", updateTarget)
    r.DELETE("/targets/:id", deleteTarget)

    r.PUT("/complete_mission/:id", completeMission)

    r.Run(":8080")
}
