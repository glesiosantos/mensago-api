package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Product struct {
	Id          string
	Description string
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		obj := map[string]string{"message": "success"}
		render.JSON(w, r, obj)
	})

	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		var product Product

		product.Id = "1"
		render.DecodeJSON(r.Body, &product)
		render.JSON(w, r, product)
	})

	//param
	// http://localhost:3000/1
	r.Get("/{idCampaign}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "idCampaign")
		w.Write([]byte(param))
	})

	//Query
	// http://localhost:3000/query?campaign=teste
	r.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("campaign")
		w.Write([]byte(query))
	})

	http.ListenAndServe(":3000", r)
}
