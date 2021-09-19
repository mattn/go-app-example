package main

import (
	"fmt"

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
