package test

const dbType = "postgres"
const dbVersion = "11.4"
const dbPortMapping = "5432/tcp"

const dbUser = "test_forna_user"
const dbUserVar = "POSTGRES_USER=" + dbUser

const dbPassword = "12345"
const dbPasswordVar = "POSTGRES_PASSWORD=" + dbPassword

const dbName = "test_forna"
const dbNameVar = "POSTGRES_DB=" + dbName

const dbURLTemplate = dbType + "://" + dbUser + ":" + dbPassword + "@localhost:%s/" + dbName + "?sslmode=disable"
