package database

import "fmt"

type DBError struct { // ? should we make separated errors such as migration error, write error, read error and ...
	TableName string
	Reason    string
}

func (e DBError) Error() string {
	return fmt.Sprintf("database error occurred on db path %s table %s with error: %s",
		e.TableName, e.Reason)
}
