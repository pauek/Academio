
package main

import (
	"Academio/webapp/data"
	"fmt"
)

func main() {
	var login, passwd string
	fmt.Printf("login: ")
	fmt.Scanf("%s", &login)
	fmt.Printf("passwd: ")
	fmt.Scanf("%s", &passwd)
	_, err := data.AddUser(login, passwd)
	if err != nil {
		if err == data.ErrUserExists {
			if user := data.GetUser(login); user != nil {
				fmt.Printf("User exists: %s [%s]\n", user.Login, user.Hpasswd)
			}
		} else {
			fmt.Printf("error: %s\n", err)
		}
	}
}