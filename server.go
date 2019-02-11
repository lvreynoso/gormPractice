// server.go

package main

import (
"time"
"errors"
"fmt"
"encoding/json"
"io/ioutil"
"log"
"net/http"
"github.com/labstack/echo"
"strings"
"math/rand"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/sqlite"
)
    
// STRETCH CHALLENGE OPTIONS:
// [DONE] return an array of words in a ChuckJoke
// [DONE] call another API with similar content to original ChuckJoke

// Takes in Chuck Norris API; used in Taco struct
type ChuckJoke struct {
    ID int `json:"id"`
    Joke string `json:"joke"`
    Categories []string `json:"categories"`
}

// Takes in ChuckJoke struct; used in texasRanger
type Taco struct {
    Type string `json:"type"`
    Value ChuckJoke `json:"value"`
}

// Takes in tronalddump API; used in TrumpDump struct
type Dumps struct {
    Value string `json:"value"`
}

// Takes in TrumpDump struct; used in func newYorkBarFly
type TrumpQuotes struct {
    Embedded struct {
        Quotes []Dumps `json:"quotes"`
    } `json:"_embedded"`
}

func texasRanger() string {
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

    return taco.Value.Joke
}

func newYorkBarFly(vomit []string) (_ string, err error) {

    fmt.Println("WELCOME TO NEW YORK")

    chalupa := TrumpQuotes{}

    for index, word := range vomit {

        fmt.Printf("FOR THE %d TIME\n", index + 1)

        if len(word) < 4 {
            fmt.Println("Wow, this word sucks! It has 3 or fewer characters in it!")
            continue
        }

        fmt.Println("Our word is longer than 3 characters!")

        response, responseErr := http.Get("https://api.tronalddump.io/search/quote?query=" + word)
        if responseErr != nil {
            fmt.Println("Error getting a response")
            log.Fatalln(err)
            continue
        }

        fmt.Println("We got a response!")

        body, bodyErr := ioutil.ReadAll(response.Body)
        if bodyErr != nil {
            fmt.Println("Error reading the response body")
            log.Fatalln(err)
            continue
        }

        fmt.Println("We read the body!")

        rawChalupa := TrumpQuotes{}

        jsonErr := json.Unmarshal(body, &rawChalupa)
        if jsonErr != nil {
            fmt.Println("Oh shit! JSON couldn't be unmarsheled! Trying again!")
            continue
        } else if len(rawChalupa.Embedded.Quotes) < 1 {
            fmt.Println("No quotes! Try again!")
            continue
        } else {
            fmt.Println("We unmarsheled that JSON and got some quotes!")
            chalupa = rawChalupa
            break
        }
    }

    fmt.Println("Chalupa:")
    fmt.Println(chalupa.Embedded.Quotes)

    trumpQuote := ""

    if len(chalupa.Embedded.Quotes) > 0 {
        err = nil
        trumpQuote = chalupa.Embedded.Quotes[rand.Intn(len(chalupa.Embedded.Quotes))].Value
    } else {
        fmt.Println("Wow this quote sucks! Tronald Dump hasn't said anything about any of these words!")
        err = errors.New("Chuck Norris? Never heard of him.")
    }

    return trumpQuote, err
}

func main() {
    server := echo.New()

    rand.Seed(time.Now().Unix())

    server.GET("/", func(context echo.Context) error {
        tacoFilling := texasRanger()
        groundBeef := strings.Split(tacoFilling, " ")
        meatMap := make(map[string]int)
        for _, beef := range groundBeef {
            if _, ok := meatMap[beef]; ok {
                meatMap[beef] += 1
            } else {
                meatMap[beef] = 1
            }
        }

        hotWord := make([]string, 0)
        for key, value := range meatMap {
            if value == 1 {
                hotWord = append(hotWord, key)
            }
        }
        fmt.Println(`Chuck's words are:`)
        fmt.Println(hotWord)
        upchuck, err := newYorkBarFly(hotWord)

        if err != nil {
            upchuck = err.Error()
        }

        return context.HTML(http.StatusOK, "<p><em>" + tacoFilling + "</em>" + `<br><em>- Faith Chikwekwe</em></p><br><p><em>` + upchuck + `</em><br><em>- Lucia Reynoso</em></p>`)

    })

    server.Logger.Fatal(server.Start(":9001"))
}
