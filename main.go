package main

import (
	"github.com/wehard/hive-master/game"
	"github.com/wehard/hive-master/ui"
)

func main() {
	/*
		clientCredentials := ftapi.Authorize()
		fmt.Println("Welcome, ", ftapi.GetAuthorizedUserData(clientCredentials.AccessToken).Displayname)
		campusUsers := ftapi.RequestCampusUsers(13, clientCredentials.AccessToken)

		for i := range campusUsers {
			fmt.Println(campusUsers[i].Login)
		}*/

	//level := game.LoadLevelFromFile("maps/level1.map")
	ui := &ui.UI2d{}
	game.Run(ui)
}
