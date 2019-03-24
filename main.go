package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
)

var titleCaseCache = cache.New(60*time.Minute, 1440*time.Minute)
var camelCaseCache = cache.New(60*time.Minute, 1440*time.Minute)

func ConvertToCamelCase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input = strings.ToLower(ps.ByName("word"))
	var possibleCamelCase = tryGetCamelCaseFromCache(input)
	if possibleCamelCase != "" {
		fmt.Fprintf(w, possibleCamelCase)
	} else {
		var firstWord = true
		var result = ""
		var startIndex = 0
		var currLen = (len(input))

		for {
			var currWord = input[startIndex:currLen]
			var possibleTitleCase = tryGetTitleCaseFromCache(currWord)
			if possibleTitleCase != "" {
				if firstWord {
					firstWord = false
					result += currWord
				} else {
					result += possibleTitleCase
				}
				startIndex = currLen
				currLen = (len(input))
			} else if doesWordExistInDictionary(currWord) {
				var equivalentTitleCase = strings.Title(currWord)
				titleCaseCache.Set(currWord, equivalentTitleCase, cache.DefaultExpiration)
				if firstWord {
					firstWord = false
					result += currWord
				} else {
					result += equivalentTitleCase
				}
				startIndex = currLen
				currLen = (len(input))
			} else {
				currLen = currLen - 1
			}

			if startIndex >= (len(input)) {
				break
			}
		}

		if result != "" {
			camelCaseCache.Set(input, result, cache.DefaultExpiration)
		}

		fmt.Fprintf(w, result)
	}
}

func tryGetTitleCaseFromCache(word string) string {
	equivalentTitleCase, found := titleCaseCache.Get(word)
	if found {
		return equivalentTitleCase.(string)
	}
	return ""
}

func tryGetCamelCaseFromCache(word string) string {
	equivalentCamelCase, found := camelCaseCache.Get(word)
	if found {
		return equivalentCamelCase.(string)
	}
	return ""
}

func doesWordExistInDictionary(word string) bool {
	appID := "<enter app id>"
	appKey := "<key>"

	language := "en"
	url := "https://od-api.oxforddictionaries.com:443/api/v1/entries/" + language + "/" + word

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("app_id", appID)
	request.Header.Set("app_key", appKey)

	client := &http.Client{}
	response, _ := client.Do(request)

	if response.StatusCode != 404 {
		return true
	}

	return false
}

func main() {
	router := httprouter.New()
	router.GET("/:word", ConvertToCamelCase)

	// Get port number from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}

	fmt.Println("Listen serve for port... ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	fmt.Println("Server started .....")
}
