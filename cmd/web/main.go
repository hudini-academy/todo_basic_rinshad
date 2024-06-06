package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/kpmohammedrinshad/alex_web_app/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application,wide dependencies for the
// web application. For now we'll only include fields for the two custom logger
// we'll add more to it as the build progresses
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	todos    *mysql.TodoModel
	users    *mysql.UserModel
}

func main() {

	// Define a new command,line flag with the name 'addr', a default value of
	// and some short help text explaining what the flag controls. The value of
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command,line flag for the MySQL DSN string
	dsn := flag.String("dsn", "root:root@/todos?parseTime=true", "MYSQL database")

	// Importantly, we use the flag.Parse() function to parse the command,line
	// This reads in the command,line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any error
	// encountered during parsing the application will be terminated.

	// Define a new command,line flag for the session secret (a random key whic
	// will be used to encrypt and authenticate session cookies). It should be
	// bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "secret key")

	flag.Parse()

	//creating a file for information message with specific flags and permissions or
	//it already created use the file with updation
	//here we use append to come the details in the end of the file
	f, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Use log.New() to create a logger for writing information messages. This
	// three parameters: the destination to write the logs to (os.Stdout), a st
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the fl
	// are joined using the bitwise OR operator |.
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	//creating a file for error message with specific flags and permissions or
	//it already created use the file with updation
	//here we use append to come the details in the end of the file
	el, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer el.Close()

	// Create a logger for writing error messages in the same way, but use stde
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number
	errorLog := log.New(el, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command,line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
		log.Println(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		todos:    &mysql.TodoModel{DB: db},
		users:    &mysql.UserModel{DB: db},
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logge
	// the event of any problems
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call the new app.routes() method
	}
	infoLog.Printf("starting server on new %s", *addr)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
