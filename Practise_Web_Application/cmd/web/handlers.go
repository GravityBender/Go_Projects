package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"net/http"
	"strconv"

	"snippetbox/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.html", &templateData{
		Snippets: snippets,
	})
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", &templateData{
		Snippet: snippet,
	})
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "create.html", &templateData{
		Snippets: nil,
		Form:     snippetCreateForm{Expires: 365},
	})
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If there are any validation errors re-display the create.html template,
	// passing in the snippetCreateForm instance as dynamic data in the Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	if len(form.FieldErrors) > 0 {
		data := &templateData{
			Snippet: nil,
		}
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
