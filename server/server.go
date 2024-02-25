package server

import (
	"L0/storage"
	"html/template"
	"log"
	"net/http"
)

// Запускает http сервер
func RunServer() {
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/main/data", dataHeandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}

// Обработчик главной страницы
func mainHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("server/view.html")
	if err != nil {
		log.Fatal(err)
	}
	err = html.Execute(writer, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Обрабочик страницы с заказом
func dataHeandler(writer http.ResponseWriter, request *http.Request) {

	uid := request.FormValue("uid")
	if len(storage.Cache) == 0 {
		storage.RestoreCache(storage.Cache)
	}

	data := storage.Cache[uid]
	if data == "" {
		data = "Заказ с таким id отсутсвует"
	}
	_, err := writer.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

}
