package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func validateBreed(breed string) bool {
    resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
    if err != nil {
        log.Println(err)
        return false
    }
    defer resp.Body.Close()

    var breeds []struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
        log.Println(err)
        return false
    }

    for _, b := range breeds {
        if b.Name == breed {
            return true
        }
    }

    return false
}
