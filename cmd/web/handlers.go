package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kpmohammedrinshad/alex_web_app/pkg/models"
)

// Change the signature of the home handler so it is defined as a method agains
// *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Send a template which will take input of creating task and listing all the tasks

	// Parse and Execute the template and pass TaskList as data

	s, err := app.todos.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, todo := range s {
	// 	fmt.Println(w, "%v\n", todo)
	// }

	files := []string{
		"./ui/html/list.home.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// 	// TODO: Improve layout and stying of the templates

	ts, err := template.ParseFiles(files...)
	if err != nil {

		app.serverError(w, err)
		http.Error(w, "Internal server Error", 500)
		return
	}
	err = ts.Execute(w, struct {
		Tasks []*models.Todo
		Flash string
	}{
		Tasks: s,
		Flash: app.session.PopString(r, "flash"),
	})

	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal server Error", 500)
	}
	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one,time fetch. If there is no matching key in the session
	// data this will return the empty string.

}

// Change the signature of the home handler so it is defined as a method agains
// *application
func (app *application) addTasks(w http.ResponseWriter, r *http.Request) {

	// Add a task to a list and send the list display to the frontend
	// Create some variables holding data
	name := r.FormValue("task")
	expires := "7"

	errors := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors["task"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["task"] = "this field is too long (maximum is 100 characters)"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	// Pass the data to the TodoModel.Insert() method, receiving the
	// ID of the new record back.
	_, err := app.todos.Insert(name, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	app.session.Put(r, "flash", "todo successfully created!")

	// Optional : Send a redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteTasks(w http.ResponseWriter, r *http.Request) {
	// Remove the task from the list and send list display to the frontend

	//converting string to integer
	id, _ := strconv.Atoi(r.FormValue("id"))

	err := app.todos.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "todo deleted successfully created!")

	// Optional: Send a redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) updateTasks(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("update")

	errors := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors["update"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["update"] = "this field is too long (maximum is 100 characters)"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}
	err := app.todos.Update(id, name)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "todo updated successfully created!")

	// Optional : Send a redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	// Create or initialize your template with ParseFiles or ParseGlob
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		// Handle error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template directly without passing any data
	err = tmpl.Execute(w, nil)
	if err != nil {
		// Handle error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
