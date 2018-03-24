package web

import (
	"gosiah/catalog"
	"gosiah/solr"
	"log"
	"net/http"
	"strings"
)

func home(values RouteValues, resp http.ResponseWriter, req *http.Request) {
	s := NewSession(values, resp, req)
	renderTemplate(s, "views/index.html", nil)
}

func search(values RouteValues, resp http.ResponseWriter, req *http.Request) {
	options := map[string]string{
		"defType": "edismax",
		"qf":      "authorsAll title^100",
	}
	facet := solr.FacetField{Name: "subjects_str", Title: "Subjects"}
	params := solr.SearchParams{
		Q:       strings.Join(req.URL.Query()["q"], " "),
		Rows:    20,
		Start:   0,
		Facets:  []solr.FacetField{facet},
		Options: options,
	}

	url := "http://localhost:8983/solr/bibdata"
	cat := catalog.New(url)
	results, err := cat.Search(params)
	log.Printf("Found: %d", results.NumFound)

	s := NewSession(values, resp, req)
	if err != nil {
		renderError(s, "Error during search", err)
	} else {
		renderTemplate(s, "views/results.html", results)
	}
}

func viewOne(values RouteValues, resp http.ResponseWriter, req *http.Request) {
	bib := values["bib"]
	log.Printf("fetching bib: %s", bib)

	url := "http://localhost:8983/solr/bibdata"
	cat := catalog.New(url)
	record, err := cat.Get(bib)

	s := NewSession(values, resp, req)
	if err != nil {
		renderError(s, "Error retrieving document from Solr", err)
	} else {
		renderTemplate(s, "views/one.html", record)
	}
}
