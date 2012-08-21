package data

import (
	"database/sql"
	"log"
	"strings"
	_ "github.com/bmizerany/pq"
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
)

type User struct {
	Login   string
	Hpasswd string
}

var ErrUserExists = errors.New("User already exists")

var db *sql.DB

var args = []string{
	"host=/var/run/postgresql",
	"user=pauek",
	"dbname=academio",
	"sslmode=disable",
}

func init() {
	_db, err := sql.Open("postgres", strings.Join(args, " "))
	if err != nil {
		log.Fatalf("Open: %s", err)
	}
	db = _db

	stmtGetUser, err = db.Prepare("SELECT * FROM users WHERE login=$1")
	if err != nil {
		log.Fatalf("Cannot prepare 'GetUser': %s", err)
	}

	stmtPutUser, err = db.Prepare("INSERT INTO users VALUES ($1, $2)")
	if err != nil {
		log.Fatalf("Cannot prepare 'PutUser': %s", err)
	}
}

var stmtGetUser *sql.Stmt

func GetUser(login string) *User {
	rows, err := stmtGetUser.Query(login)
	rows.Next()
	var user User
	if err = rows.Scan(&user.Login, &user.Hpasswd); err != nil {
		return nil
	}
	return &user
}

var stmtPutUser *sql.Stmt

func PutUser(user *User) error {
	_, err := stmtPutUser.Exec(user.Login, user.Hpasswd)
	if err != nil {
		return err
	}
	return nil
}

func AuthenticateUser(login, password string) *User {
	if user := GetUser(login); user != nil {
		a := []byte(user.Hpasswd)
		b := []byte(password)
		if err := bcrypt.CompareHashAndPassword(a, b); err == nil {
			return user
		}
	}
	return nil
}

func AddUser(login, pass string) (user *User, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}
	if existingUser := GetUser(login); existingUser != nil {
		return nil, ErrUserExists
	}
	user = &User{
		Login:   login,
		Hpasswd: string(hash),
	}
	if err := PutUser(user); err != nil {
		return nil, err
	}
	return user, nil
}
