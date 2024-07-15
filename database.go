package main

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "os"
    "fmt"
    "github.com/joho/godotenv"
)

var db *sqlx.DB

func initDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Помилка завантаження .env файлу: %v", err)
    }

    user := os.Getenv("USER")
    password := os.Getenv("PASSWORD")
    dbname := os.Getenv("DBNAME")
    host := os.Getenv("HOST")
    port := os.Getenv("PORT")
    sslmode := os.Getenv("SSLMODE")

    connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
        user, password, dbname, host, port, sslmode)

    db, err = sqlx.Connect("postgres", connectionString)
    if err != nil {
        log.Fatalln(err)
    }

    schema := `
    CREATE TABLE IF NOT EXISTS cats (
        id SERIAL PRIMARY KEY,
        name TEXT,
        experience INTEGER,
        breed TEXT,
        salary REAL
    );

    CREATE TABLE IF NOT EXISTS missions (
        id SERIAL PRIMARY KEY,
        cat_id INTEGER REFERENCES cats(id),
        complete BOOLEAN
    );

    CREATE TABLE IF NOT EXISTS targets (
        id SERIAL PRIMARY KEY,
        mission_id INTEGER REFERENCES missions(id),
        name TEXT,
        country TEXT,
        notes TEXT,
        complete BOOLEAN
    );
    `
    
    db.MustExec(schema)
}
