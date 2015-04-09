package main

import (
	"fmt"
	"github.com/jothirams/go-alfred"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

	// Normalize the queryTerm - because über comes inside as u¨ber
	// Refer: http://alfredworkflow.readthedocs.org/en/latest/user-manual/text-encoding.html
	// http://blog.golang.org/normalization
	verbs, err := getVerbList(norm.NFC.String(queryTerm))

	if err != nil {
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: false,
			Title: err.Error(),
		})
	} else {
		for _, verb := range verbs {

			// it matched so add a new response item
			response.AddItem(&alfred.AlfredResponseItem{
				Valid: true,
				Uid:   verb.URL,
				Title: verb.Name,
				Arg:   verb.URL,
			})
		}
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
func getVerbList(queryTerm string) ([]VerbList, error) {

	// Encode URL with queryTerm
	resp, err := http.Get(fmt.Sprint("http://www.verblisten.de/eingabeliste.jsp?eingabe=", url.QueryEscape(queryTerm)))
	if err != nil {
		return nil, fmt.Errorf("Unable to reach verblisten.de. Failed with error: %s\"", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("The received response body is has returned error: %s\"", err.Error())
	}

	sbody := string(body)
	if len(sbody) <= 0 {
		return nil, fmt.Errorf("No matching verbs for \"%s\"", queryTerm)
	}

	regExp := regexp.MustCompile("[äöüß]")

	umlauts := map[string]string{
		"ä": "a:",
		"ö": "o:",
		"ü": "u:",
		"ß": "s:",
	}

	replaceUmlauts := func(str string) string {
		return umlauts[str]
	}

	verbs := make([]VerbList, 10, 10)

	j := 0
	for _, s := range strings.Split(sbody, ";") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			verbs[j].Name = s
			urlQueryTerm := regExp.ReplaceAllStringFunc(s, replaceUmlauts)
			verbs[j].URL = fmt.Sprintf("http://www.verbformen.de/konjugation/%s.htm", urlQueryTerm)
			j++
			if j == 10 {
				break
			}
		}
	}

	return verbs, nil
}
