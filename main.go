package main

import (
	"encoding/json"
	"fmt"

	"github.com/wehard/ftapi"
)

func main() {

	clientCredentials := ftapi.Authorize()
	fmt.Println("token:", clientCredentials.AccessToken)
	bytes := ftapi.DoFTRequest("/v2/me", clientCredentials.AccessToken)
	var userData ftapi.UserData
	json.Unmarshal(bytes, &userData)
	fmt.Println(userData.Displayname)
}
