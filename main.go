package main

import (
	"fmt"

	"github.com/wehard/ftapi"
)

func main() {
	clientCredentials := ftapi.Authorize()
	fmt.Println("Welcome, ", ftapi.GetAuthorizedUserData(clientCredentials.AccessToken).Displayname)
	ftapi.GetCampusUsers(13, clientCredentials.AccessToken)
}
