package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type item struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type todoList struct {
	app.Compo

	items []item
}

func newTodoList() *todoList {
	return &todoList{}
}

func (l *todoList) OnMount(ctx app.Context) {
	ctx.Async(func() {
		resp, err := http.Get("/api")
		if err != nil {
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&l.items)
		if err != nil {
			return
		}
		l.Update()
	})
}

func (h *todoList) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Todo List"),
		),
		app.Input().
			Value("").
			OnChange(h.OnInputChange),
		app.P().
			Class("center-content").
			Body(
				app.Range(h.items).Slice(func(i int) app.UI {
					v := h.items[i]
					return app.Div().Class("todo-item").Body(
						app.Input().
							ID(fmt.Sprintf("item-%d", v.ID)).
							Type("checkbox").
							Checked(v.Done).
							Class("button").
							OnChange(h.OnDoneChange),
						app.Div().
							Text(v.Text),
					)
				}),
			),
	)
}

func (h *todoList) OnInputChange(ctx app.Context, e app.Event) {
	text := ctx.JSSrc.Get("value").String()

	b, err := json.Marshal(item{Text: text})
	if err != nil {
		return
	}

	resp, err := http.Post("/api", "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	h.items = append(h.items, item{ID: len(h.items) + 1, Text: text})
	ctx.JSSrc.Set("value", "")
	h.Update()
}

func (h *todoList) OnDoneChange(ctx app.Context, e app.Event) {
	s := ctx.JSSrc.Get("id").String()
	if !strings.HasPrefix(s, "item-") {
		return
	}
	id, err := strconv.Atoi(s[5:])
	if err != nil {
		return
	}
	done := ctx.JSSrc.Get("checked").Bool()

	b, err := json.Marshal(item{ID: id, Done: done})
	if err != nil {
		return
	}

	resp, err := http.Post("/api", "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func listTodo(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("todo.json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	var items []item
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	var v item
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	f, err := os.Open("todo.json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var items []item
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(f).Encode(&items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	f.Close()
}

func main() {
	app.Route("/", &todoList{})

	app.RunWhenOnBrowser()

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
