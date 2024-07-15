package main

import (

    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

type Target struct {
    ID       int    `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
    Country  string `json:"country" db:"country"`
    Notes    string `json:"notes" db:"notes"`
    Complete bool   `json:"complete" db:"complete"`
}

type Mission struct {
    ID       int      `json:"id" db:"id"`
    CatID    int      `json:"cat_id" db:"cat_id"`
    Complete bool     `json:"complete" db:"complete"`
    Targets  []Target `json:"targets"`
}

func createMission(c *gin.Context) {
    var mission Mission
    if err := c.ShouldBindJSON(&mission); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tx := db.MustBegin()

    query := `INSERT INTO missions (cat_id, complete) VALUES ($1, $2) RETURNING id`
    err := tx.QueryRow(query, mission.CatID, false).Scan(&mission.ID)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    newTargets := []Target{}

    for _, target := range mission.Targets {
        query = `INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5) RETURNING id`
        err = tx.QueryRow(query, mission.ID, target.Name, target.Country, target.Notes, false).Scan(&target.ID)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        newTargets = append(newTargets, target)
    }

    mission.Targets = newTargets

    err = tx.Commit()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, mission)
}

func getMissions(c *gin.Context) {
    missions := []Mission{}
    err := db.Select(&missions, "SELECT id, cat_id, complete FROM missions")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    for i, mission := range missions {
        err := db.Select(&missions[i].Targets, "SELECT id, name, country, notes, complete FROM targets WHERE mission_id = $1", mission.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    c.JSON(http.StatusOK, missions)
}

func getMission(c *gin.Context) {
    id := c.Param("id")
    var mission Mission
    query := `SELECT id, cat_id, complete FROM missions WHERE id = $1`
    err := db.Get(&mission, query, id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
        return
    }

    err = db.Select(&mission.Targets, "SELECT id, name, country, notes, complete FROM targets WHERE mission_id = $1", mission.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, mission)
}

func updateMission(c *gin.Context) {
    id := c.Param("id")
    var mission Mission
    if err := c.ShouldBindJSON(&mission); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE missions SET complete = $1 WHERE id = $2`
    _, err := db.Exec(query, mission.Complete, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func deleteMission(c *gin.Context) {
    id := c.Param("id")

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    deleteTargetsQuery := `DELETE FROM targets WHERE mission_id = $1`
    _, err = tx.Exec(deleteTargetsQuery, id)
    if err != nil {
        tx.Rollback()
        log.Printf("Error deleting targets: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    deleteMissionQuery := `DELETE FROM missions WHERE id = $1`
    _, err = tx.Exec(deleteMissionQuery, id)
    if err != nil {
        tx.Rollback()
        log.Printf("Error deleting mission: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Error committing transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Printf("Mission with ID: %s successfully deleted", id)
    c.Status(http.StatusNoContent)
}

func addTarget(c *gin.Context) {
    missionID := c.Param("mission_id")
    var target Target
    if err := c.ShouldBindJSON(&target); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    err := db.QueryRow(query, missionID, target.Name, target.Country, target.Notes, false).Scan(&target.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, target)
}




func completeMission(c *gin.Context) {
    id := c.Param("id")

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    completeTargetsQuery := `UPDATE targets SET complete = true WHERE mission_id = $1`
    _, err = tx.Exec(completeTargetsQuery, id)
    if err != nil {
        tx.Rollback()
        log.Printf("Error completing targets: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    completeMissionQuery := `UPDATE missions SET complete = true WHERE id = $1`
    _, err = tx.Exec(completeMissionQuery, id)
    if err != nil {
        tx.Rollback()
        log.Printf("Error completing mission: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Error committing transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Printf("Mission with ID: %s successfully completed", id)
    c.Status(http.StatusNoContent)
}

func updateTarget(c *gin.Context) {
    id := c.Param("id")
    var target Target
    if err := c.ShouldBindJSON(&target); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var complete bool
    query := `SELECT complete FROM targets WHERE id = $1`
    err := db.Get(&complete, query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if complete {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update completed target"})
        return
    }

    query = `UPDATE targets SET notes = $1, complete = $2 WHERE id = $3`
    _, err = db.Exec(query, target.Notes, target.Complete, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    missionID := 0
    query = `SELECT mission_id FROM targets WHERE id = $1`
    err = db.Get(&missionID, query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var count int
    query = `SELECT COUNT(*) FROM targets WHERE mission_id = $1 AND complete = false`
    err = db.Get(&count, query, missionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if count == 0 {
        query = `UPDATE missions SET complete = true WHERE id = $1`
        _, err = db.Exec(query, missionID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    c.Status(http.StatusNoContent)
}


func deleteTarget(c *gin.Context) {
    id := c.Param("id")
    query := `DELETE FROM targets WHERE id = $1`
    _, err := db.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
