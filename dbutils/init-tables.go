package dbutils
import  "log"
import "database/sql"

func Initialize(dbDriver *sql.DB) {
    statement, driverError := dbDriver.Prepare(user)
    if driverError != nil {
        log.Println(driverError)
    }
    // Create user table
    _, statementError := statement.Exec()
    if statementError != nil {
        log.Println("User Table already exists!")
    }

}
