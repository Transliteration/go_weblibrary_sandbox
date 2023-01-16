package database

import (
	"database/sql"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	pb "weblibrary_sandbox/grpc_server"

	"github.com/google/uuid"
	// _ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	_ "github.com/go-sql-driver/mysql"
)

const (
	DBFilename = "sqlite-database.db"
)

var db *sql.DB

// sqlite init
func createDBFileIfItNotExist() (bool, error) {
	info, err := os.Stat(DBFilename)
	if os.IsNotExist(err) {
		log.Print("DB file doesnt exist, creating new...")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Print("File created!")

		return true, nil

	} else if info.IsDir() {
		return false, fmt.Errorf("%s is directory, not db file", DBFilename)
	}

	return false, nil
}

func InitDB() error {

	// DBjustCreated, err := createDBFileIfItNotExist()
	// if err != nil {
	// 	return err
	// }

	var err error
	db, err = sql.Open("mysql", "user:password@/grpctest_database")
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	populateDBOnCreation()

	fmt.Println("Connection to DB established!")

	return nil
}

func CloseDB() {
	db.Close()
}

func populateDBOnCreation() {
	createTable()

	row, err := db.Query("SELECT COUNT(*) FROM userInfo;")
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer row.Close()

	var count int
	row.Next()
	row.Scan(&count)

	if count == 0 {
		AddUser("Max", 44)
		AddUser("Vasya", 22)
		AddUser("Nikolay", 34)
		AddUser("Oleg", 51)
		AddUser("Flex", 13)
	}

	displayUsers()
}

func createTable() error {
	log.Println("Creating users table...")
	createTableUsers := `CREATE TABLE IF NOT EXISTS userInfo (
		userID VARCHAR(255) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		age INT
	);`

	// createTableUsers = `DROP TABLE userInfo;`

	statement, err := db.Prepare(createTableUsers)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer statement.Close()

	statement.Exec()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Println("Users table created!")

	return nil
}

func AddUser(name string, age int) (*pb.UserId, error) {
	log.Println("Inserting user record ...")

	insertUserStatement := `INSERT INTO userInfo(userID, name, age) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertUserStatement)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer statement.Close()

	newUserID := uuid.New().String()
	_, err = statement.Exec(newUserID, name, age)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &pb.UserId{Id: newUserID}, nil
}

func GetAllUsers() (*pb.Users, error) {
	var allUsers []*pb.User

	row, err := db.Query("SELECT * FROM userInfo ORDER BY name")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer row.Close()

	for row.Next() { // Iterate and fetch the records from result cursor
		var userID string
		var name string
		var age int

		err = row.Scan(&userID, &name, &age)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		user := &pb.User{UserId: userID, Name: name, Age: int32(age)}
		allUsers = append(allUsers, user)
	}

	return &pb.Users{Users: allUsers}, nil
}

func GetUser(id string) (*pb.User, error) {
	stm, err := db.Prepare("SELECT * FROM userInfo WHERE userID = ?")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stm.Close()

	user := &pb.User{}
	err = stm.QueryRow(id).Scan(&user.UserId, &user.Name, &user.Age)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}

func displayUsers() error {
	row, err := db.Query("SELECT * FROM userInfo ORDER BY name")
	if err != nil {
		log.Error(err)
		return err
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var userID string
		var name string
		var age int
		err = row.Scan(&userID, &name, &age)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Println("User: ", name, " age = ", age, " id = ", userID)
	}

	return nil
}
