//go:build wasm

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

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

func main() {
	app.Route("/", &todoList{})
	app.RunWhenOnBrowser()
}
