//go:build !wasm

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

func (l *todoList) OnMount(ctx app.Context) {
}

func httpError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func listTodo(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("todo.json")
	if err != nil {
		httpError(w)
		return
	}
	defer f.Close()

	var items []item
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		httpError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&items)
	if err != nil {
		httpError(w)
		return
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	var v item
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		httpError(w)
		return
	}

	f, err := os.Open("todo.json")
	if err != nil {
		httpError(w)
		return
	}
	var items []item
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		httpError(w)
		return
	}
	f.Close()

	found := false
	for i, vv := range items {
		if vv.ID == v.ID {
			items[i].Done = v.Done
			found = true
			break
		}
	}
	if !found {
		items = append(items, item{ID: len(items) + 1, Text: v.Text})
	}
	f, err = os.Create("todo.json")
	if err != nil {
		httpError(w)
		return
	}
	err = json.NewEncoder(f).Encode(&items)
	if err != nil {
		httpError(w)
		return
	}
	f.Close()
}

func (h *todoList) OnInputChange(ctx app.Context, e app.Event) {
}

func (h *todoList) OnDoneChange(ctx app.Context, e app.Event) {
}

func main() {
	app.Route("/", &todoList{})

	http.Handle("/", &app.Handler{
		Name:        "Todo List",
		Description: "Todo List",
		Styles: []string{
			"/web/style.css",
		},
	})

	var mu sync.Mutex
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			mu.Lock()
			defer mu.Unlock()
			listTodo(w, r)
		} else if r.Method == http.MethodPost {
			mu.Lock()
			defer mu.Unlock()
			updateTodo(w, r)
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
