// server.go

package main

import (
"fmt"
"encoding/json"
"io/ioutil"
"log"
"net/http"
"github.com/labstack/echo"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/sqlite"
)
    
// STRETCH CHALLENGE OPTIONS:
// [DONE] return an array of words in a ChuckJoke
// [DONE] call another API with similar content to original ChuckJoke

// Takes in Chuck Norris API; used in Taco struct
type Norris struct {
    gorm.Model
    JokeID int `json:"id" form:"jokeid" query:"jokeid"`
    Joke string `json:"joke" form:"joke" query:"joke"`
    // Categories []string `json:"categories"`
}

// Takes in ChuckJoke struct; used in texasRanger
type Taco struct {
    Type string `json:"type"`
    Value Norris `json:"value"`
}

func texasRanger() Taco {
    // takes in Taco struct and returns a Chuck Norris jok as a string
    response, err := http.Get("https://api.icndb.com/jokes/random")
    if err != nil {
        log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatalln(err)
    }

    taco := Taco{}

    jsonErr := json.Unmarshal(body, &taco)
    if jsonErr != nil {
        log.Fatalln(jsonErr)
    }

    return taco
}

func main() {
    server := echo.New()

    db, err := gorm.Open("sqlite3", "practice.db")
      if err != nil {
        panic("failed to connect database")
      }
    defer db.Close()

    db.AutoMigrate(&Norris{})

    server.GET("/", func(context echo.Context) error {
        rows, err := db.Model(&Norris{}).Rows()
        defer rows.Close()
        if err != nil {
            fmt.Println(err)
            return context.HTML(http.StatusInternalServerError, "Womp womp")
        }

        jokeHTML := "<p>"

        for rows.Next() {
            var norris Norris
            db.ScanRows(rows, &norris)
            jokeHTML += norris.Joke + "<br>"
        }

        jokeHTML += "</p>"

        return context.HTML(http.StatusOK, "<h1>Success!</h1><br><h4>Current Jokes:</h4>" + jokeHTML)
        })

    server.GET("/upload", func(context echo.Context) error {
            uploadForm := `
            <h3>Chuck Norris joke uploader for Golang</h3>
            <br>
            <form action="/upload" method="post" enctype="multipart/form-data">
                <div>
                    <label for="jokeid">Enter a Joke ID:</label><br>
                    <input type="text" id="jokeid" name="jokeid" value="">
                </div>
                <br>
                <div>
                    <label for="joke">Enter a Chuck Norris Joke:</label><br>
                    <input type="text" id="Joke" name="joke" value="">
                </div>
                <br>
                <div>
                    <button type="submit" name="button">Upload</button>
                </div>
            </form>`
            return context.HTML(http.StatusOK, uploadForm)
        })

    server.POST("/upload", func(context echo.Context) error {
        norris := new(Norris)
        if err = context.Bind(norris); err != nil {
            fmt.Println("Womp womped on the bind bind")
            return context.HTML(http.StatusInternalServerError, "Womp womp")
        }

        // put the joke in the table
        db.Create(&Norris{JokeID: norris.JokeID, Joke: norris.Joke})

        return context.JSON(http.StatusOK, norris)
        })

    server.GET("/update", func(context echo.Context) error {
            uploadForm := `
            <h3>Chuck Norris joke uPdAtEr for Golang</h3>
            <br>
            <form action="/update" method="PUT" enctype="multipart/form-data">
                <div>
                    <label for="jokeid">Enter the Joke ID:</label><br>
                    <input type="text" id="jokeid" name="jokeid" value="">
                </div>
                <br>
                <div>
                    <label for="joke">Enter the new Chuck Norris Joke:</label><br>
                    <input type="text" id="Joke" name="joke" value="">
                </div>
                <br>
                <div>
                    <button type="submit" name="button">Upload</button>
                </div>
            </form>`
            return context.HTML(http.StatusOK, uploadForm)
        })

    server.PUT("/update", func(context echo.Context) error {
        updated := new(Norris)
        if err = context.Bind(updated); err != nil {
            fmt.Println("Womp womped on the bind bind")
            return context.HTML(http.StatusInternalServerError, "Womp womp")
        }

        var norris Norris

        db.Where(&Norris{JokeID: updated.JokeID}).First(&norris)

        norris.Joke = updated.Joke

        db.Save(&norris)

        return context.JSON(http.StatusOK, norris)
        })

    server.GET("/populate", func(context echo.Context) error {
        taco := texasRanger()
        tacoFilling := taco.Value

        // create a new joke in SQL
        db.Create(&Norris{JokeID: tacoFilling.JokeID, Joke: tacoFilling.Joke})

        return context.HTML(http.StatusOK, "<h1>SQLized!</h1><p><em>" + tacoFilling.Joke + "</em>" + `<br><em>- Faith Chikwekwe</em></p>`)
        })

    server.Logger.Fatal(server.Start(":9001"))
}
