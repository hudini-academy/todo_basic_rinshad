package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kpmohammedrinshad/alex_web_app/pkg/forms"
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

	flash := app.session.PopString(r, "flash")

	// Pass the flash message to the template.
	app.render(w, r, "list.page.tmpl", &templateData{
		Flash: flash,
		Todos: s,
	})
}

// Change the signature of the home handler so it is defined as a method agains
// *application
func (app *application) addTasks(w http.ResponseWriter, r *http.Request) {

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to
	// a 400 Bad Request response to the user.
	error := r.ParseForm()
	if error != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("task")
	form.MaxLength("task", 100)
	//form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "list.page.tmpl", &templateData{Form: form})
		return
	}
	// Use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.
	name := r.PostForm.Get("task")
	//expires := "7"

	isSpecial := strings.Contains(name, "Special:")

	if isSpecial == true {
		_, err := app.todos.Insert(form.Get("task"))
		if err != nil {
			app.serverError(w, err)
			return
		}
		_, err = app.specials.InsertSpecialTask(form.Get("task"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	} else {
		_, err := app.todos.Insert(form.Get("task"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	errors := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors["task"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["task"] = "this field is too long (maximum is 100 characters)"
	}

	if len(errors) > 0 {
		app.render(w, r, "list.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Pass the data to the TodoModel.Insert() method, receiving the
	// ID of the new record back.
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
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// Otherwise send a placeholder response (for now!).
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// Try to create a new user record in the database. If the email already exi
	// add an error message to the form and re,display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming tha
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic e
	// message to the form failures map and re,display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))

	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logg
	// in'.
	app.session.Put(r, "userID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the userID from the session data so that the user is 'logged out'
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", 303)

}

// create a function that show a special task
func (app *application) specialTask(w http.ResponseWriter, r *http.Request) {

	s, err := app.specials.LatestSpecialTask()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Pass the flash message to the template.
	flash := app.session.PopString(r, "flash")

	app.render(w, r, "special.page.tmpl", &templateData{
		Flash:    flash,
		Specials: s,
	})
}

func (app *application) specialDeleteTask(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.FormValue("id"))

	err := app.specials.DeleteSpecialTask(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "todo deleted successfully created!")

	// Optional: Send a redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
