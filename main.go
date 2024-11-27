package main

import "github.com/docker/docker-credential-helpers/credentials"

// "github.com/docker/docker-credential-helpers/credentials"

// "github.com/jfrog/jfrog-cli-core/v2/common/commands"
// "github.com/jfrog/jfrog-cli-core/v2/utils/config"

func main() {
	// var err error
	// serverId := "agilityrobotics"

	// serverDetails, err := config.GetSpecificConfig(serverId, false, false)

	// if err != nil {
	// 	serverUrl := fmt.Sprintf("https://%s.jfrog.io", serverId)
	// 	serverDetails := config.ServerDetails{Url: serverUrl}

	// 	err := commands.NewConfigCommand(commands.AddOrEdit, serverId).
	// 		SetInteractive(true).SetUseWebLogin(true).
	// 		SetDetails(&serverDetails).Run()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(serverDetails.AccessToken)

	credentials.Serve(ArtifactoryKeychain{})
}
