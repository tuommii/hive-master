package main

import (
	"fmt"
	"os"

	"github.com/wehard/ftapi"
	"github.com/wehard/hive-master/game"
	"github.com/wehard/hive-master/ui"
)

//var AuthorizedClientCredentials ftapi.ClientCredentials

func main() {

	clientCredentials := ftapi.Authorize()
	game.AuthorizedClientCredentials = clientCredentials
	fmt.Println("Welcome, ", ftapi.GetAuthorizedUserData(clientCredentials.AccessToken).Displayname)

	_, err := os.Stat("game/users.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("users file does not exist!")
			userData := ftapi.RequestAllCampusUsersData(ftapi.Hive, clientCredentials.AccessToken)
			ftapi.SaveUserData("game/users.json", userData)
		}
	}

	ui := &ui.UI2d{}
	game.Run(ui)
}
