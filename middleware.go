package main

import (
	"html/template"
	"log"
	"net/http"
)

func showPage(tmpName string, sd siteData, w http.ResponseWriter, req *http.Request) {
	t := template.New(tmpName)
	//t.Delims("{(", ")}")

	t, err := t.ParseFiles("templates/" + tmpName)
	if err != nil {
		log.Println(err.Error())
	}

	err = t.Execute(w, sd)
	if err != nil {
		log.Println(err.Error())
	}
}
