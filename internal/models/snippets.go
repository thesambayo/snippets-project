package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Snippet type is defined to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets table?
type Snippet struct {
	ID      int
	Title   string // replace with sql.NullString if column in DB can be nullable
	Content string // replace with sql.NullString if column in DB can be nullable
	Created time.Time
	Expires time.Time
}

// SnippetModel type is defined which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error instead
		// (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Latest will return the 10 most recently created  snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// the SQL statement we want to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultSet containing the result
	// of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultSet is
	// always properly closed before the Latest() method returns. This defers
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultSet.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	var snippets []*Snippet

	// Use rows.Next to iterate through the rows in the resultSet. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultSet automatically closes itself and frees-up the underlying
	// database connection.

	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}
		fmt.Println(s.Created)
		// Use rows.Scan() to copy the values from each field in the row to
		// new Snippet object that we created. Again, the arguments to row.Scan() the
		// must be pointers to the place you want to copy the data into, and
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultSet.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
