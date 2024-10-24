package main

import (
	"log"
	"net/http"
)

func main() {
	var err error
	db, err = openDatabase("./calories.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables()
	http.HandleFunc("/", renderFormPage)
	http.HandleFunc("/generate_chart", generateChartHandler)
	http.Handle("/weekly_calories.html", http.FileServer(http.Dir(".")))

	log.Println("Servidor rodando em http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}