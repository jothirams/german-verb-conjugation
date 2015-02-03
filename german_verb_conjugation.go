package main

import (
	"fmt"
	"github.com/jothirams/go-alfred"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage:", os.Args[0], "<verb_search_query>")
		os.Exit(1)
	}
	queryTerm := os.Args[1]

	// create a new alfred workflow response
	response := alfred.NewResponse()
	verbs := getVerbList(queryTerm)

	for _, verb := range verbs {

		// it matched so add a new response item
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: true,
			Uid:   verb.URL,
			Title: verb.Name,
			Arg:   verb.URL,
		})
	}

	// finally print the resulting Alfred Workflow XML
	response.Print()
}

type VerbList struct {
	Name string
	URL  string
}

// Gets the verb list from verbformen.de
// And returns with the (constructed) URL and Name
func getVerbList(queryTerm string) []VerbList {

	resp, err := http.Get(fmt.Sprint("http://www.verbformen.de/eingabeliste.jsp?eingabe=", queryTerm))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	sbody := string(body)
	verbs := make([]VerbList, 10, 10)

	j := 0
	for _, s := range strings.Split(sbody, ";") {
		if len(s) > 0 {
			verbs[j].Name = s
			verbs[j].URL = fmt.Sprintf("http://www.verbformen.de/konjugation/%s.htm", s)
			j++
			if j == 10 {
				break
			}
		}
	}

	return verbs
}