package mysql

import (
	"database/sql"
	"log"

	"github.com/kpmohammedrinshad/alex_web_app/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

// This will insert a new todo into the database.
func (t *TodoModel) Insert(name string) (int, error) {
	log.Println(name)
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	stmt := `INSERT INTO todos(name,created,expires)
	values(?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// name and expiry values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	result, err := t.DB.Exec(stmt, name, 7)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the todos table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

func (t *TodoModel) Latest() ([]*models.Todo, error) {

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, name, created, expires FROM todos
	WHERE expires > UTC_TIMESTAMP()`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := t.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Todo objects.
	todos := []*models.Todo{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees,up the underlying
	// database connection.
	for rows.Next() {

		// Create a pointer to a new zeroed Todo struct.
		s := &models.Todo{}

		// Use rows.Scan() to copy the values from each field in the row to the
		// new todo object that we created. Again, the arguments to row.Scan
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Name, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of todos.
		todos = append(todos, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this , don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the todos slice.
	return todos, nil
}

func (t *TodoModel) Delete(id int) error {

	// Write the SQL statement we want to execute.
	stmt := `DELETE FROM todos WHERE id=?`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// id value for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	_, err := t.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoModel) Update(id int, name string) error {

	// Write the SQL statement we want to execute.
	stmt := `UPDATE todos
	SET name=?
	WHERE id=?`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// name and id  values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	_, err := t.DB.Exec(stmt, name, id)
	if err != nil {
		return err
	}
	return nil
}
