package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Cat struct {
    ID         int     `json:"id" db:"id"`
    Name       string  `json:"name" db:"name"`
    Experience int     `json:"experience" db:"experience"`
    Breed      string  `json:"breed" db:"breed"`
    Salary     float64 `json:"salary" db:"salary"`
}

func createCat(c *gin.Context) {
    var cat Cat
    if err := c.ShouldBindJSON(&cat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if !validateBreed(cat.Breed) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid breed"})
        return
    }

    query := `INSERT INTO cats (name, experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id`
    err := db.QueryRow(query, cat.Name, cat.Experience, cat.Breed, cat.Salary).Scan(&cat.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, cat)
}

func getCats(c *gin.Context) {
    cats := []Cat{}
    err := db.Select(&cats, "SELECT id, name, experience, breed, salary FROM cats")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, cats)
}

func getCat(c *gin.Context) {
    id := c.Param("id")
    var cat Cat
    query := `SELECT id, name, experience, breed, salary FROM cats WHERE id = $1`
    err := db.Get(&cat, query, id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
        return
    }
    c.JSON(http.StatusOK, cat)
}

func updateCat(c *gin.Context) {
    id := c.Param("id")
    var cat Cat
    if err := c.ShouldBindJSON(&cat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE cats SET salary = $1 WHERE id = $2`
    _, err := db.Exec(query, cat.Salary, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func deleteCat(c *gin.Context) {
    id := c.Param("id")
    query := `DELETE FROM cats WHERE id = $1`
    _, err := db.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
