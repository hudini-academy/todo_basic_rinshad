package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	// Swap the route declarations to use the application struct's methods as t
	// handler functions.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Post("/addTask", dynamicMiddleware.ThenFunc(app.addTasks))
	mux.Post("/deleteTask", dynamicMiddleware.ThenFunc(app.deleteTasks))
	mux.Post("/updateTask", dynamicMiddleware.ThenFunc(app.updateTasks))

	// Add the five new routes.
	mux.Get("/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	// Create a file server which serves files out of the "./ui/static" directo
	// Note that the path given to the http.Dir function is relative to the pro
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	// Running secureHeader and then running mux
	return standardMiddleware.Then(mux)
}
