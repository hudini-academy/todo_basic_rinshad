package mysql

import (
	"database/sql"

	"github.com/kpmohammedrinshad/alex_web_app/pkg/models"
)

type SpecialModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table.
func (t *SpecialModel) InsertSpecialTask(name string) (int, error) {
	stmt := `INSERT INTO specialTask (name)
VALUES(?)`
	// Use the Exec() method to insert the user details and hashed password
	// into the users table. If this returns an error, we try to type assert
	// it to a *mysql.MySQLError object so we can check if the error number is
	// 1062 and, if it is, we also check whether or not the error relates to
	// our users_uc_email key by checking the contents of the message string.
	// If it does, we return an ErrDuplicateEmail error. Otherwise, we just
	// return the original error (or nil if everything worked).
	result, err := t.DB.Exec(stmt, name)
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

func (t *SpecialModel) LatestSpecialTask() ([]*models.SpecialTask, error) {

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, name FROM specialTask`

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
	specials := []*models.SpecialTask{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees,up the underlying
	// database connection.
	for rows.Next() {

		// Create a pointer to a new zeroed Todo struct.
		s := &models.SpecialTask{}

		// Use rows.Scan() to copy the values from each field in the row to the
		// new todo object that we created. Again, the arguments to row.Scan
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of todos.
		specials = append(specials, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this , don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the todos slice.
	return specials, nil
}

func (t *SpecialModel) DeleteSpecialTask(id int) error {

	// Write the SQL statement we want to execute.
	stmt := `DELETE FROM specialTask WHERE id=?`

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
