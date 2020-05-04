package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
)

// Package contains the fields conveying information about public packages.
type Package struct {
	Name        string
	Path        string
	ImportCount int
	Synopsis    string
	Stars       int
	Score       float32
}

// Print displays the fields of Package.
func (p *Package) Print() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight)
	println()
	fmt.Fprintln(w, "NAME:", "\t", p.Name)
	fmt.Fprintln(w, "PATH:", "\t", p.Path)
	fmt.Fprintln(w, "SYNOPSIS:", "\t", p.Synopsis)
	fmt.Fprintln(w, "STARS:", "\t", p.Stars)
	fmt.Fprintln(w, "SCORE:", "\t", p.Score)
	w.Flush()
}

// Search represents the operation of searching on the GoDoc API.
type Search struct {
	Query string
	Count int
}

// Run makes a HTTP GET request against the GoDoc API to retrive search results.
func (s *Search) Run() *Response {
	var body Response

	res, err := http.Get("https://api.godoc.org/search?q=" + url.QueryEscape(s.Query))
	if err != nil {
		log.Fatalf("Error running search against GoDoc API: %v", err)
	}

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatalf("Error decoding response from GoDoc API: %v", err)
	}

	body.Results = body.Results[:s.Count]

	return &body
}

// NewSearch instantiates a new Search instance.
func NewSearch() *Search {
	return &Search{}
}

// Response contains the fields to hold the response from GoDoc API.
type Response struct {
	Results []Package
}

func main() {
	s := NewSearch()

	// Define the parameter names, default values, and usage.
	// Write those values to the Search struct instance.
	flag.StringVar(&s.Query, "query", "", "the package search term")
	flag.StringVar(&s.Query, "q", "", "the package search term (shorthand)")
	flag.IntVar(&s.Count, "count", 6, "the number of results to return")
	flag.IntVar(&s.Count, "c", 6, "the number of results to return (shorthand)")

	flag.Parse()

	if len(os.Args) <= 1 {
		fmt.Println("Godoc-search is a program for searching packages on GoDoc.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}

	if s.Query == "" {
		log.Fatal("No query parameter `q` provided.")
	}

	body := s.Run()

	for _, pkg := range body.Results {
		pkg.Print()
	}
}
