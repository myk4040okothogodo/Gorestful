package dbutils


const user = `
    CREATE TABLE IF NOT EXISTS train (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        USER_NAME VARCHAR(64) NULL,
        EMAIL_ADDRESS VARCHAR(100) NULL,
        AGE INTEGER NULL,
        MARITAL_STATUS BOOLEAN
    )
`
