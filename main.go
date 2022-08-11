package main

import (
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/alecthomas/kong"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}


func main()  {
	defer log.Stop()

	log.Info("Bugsnag CLI Tool")

	var commands struct {
		UploadServer string `help:"Bugsnag On-Premise upload server URL"`
		ApiKey    string	`help:"Bugsnag project API key"`
		Upload       struct {
			// shared options
			Overwrite 		bool   					`help:"ignore existing upload with same version"`
			Timeout   		int    					`help:"seconds to wait before failing an upload request"`
			Retries   		int    					`help:"number of retry attempts before failing a request"`
			ExtraOptions	map[string]string		`help:"additional arguments to pass to the upload request" mapsep:","`

			// required options
			Path 	  		[]string 				`help:"Path to directory to search" arg:"" name:"path" type:"path"`
		} `cmd:"" help:"Upload files"`
	}

	ctx := kong.Parse(&commands)

	pathIsValid, path := utils.ValidatePath(commands.Upload.Path)

	//// Check if path(s) provided are valid
	if !pathIsValid {
		// If the path(s) is not valid, log it and return early
		log.Error("%s does not exist", path)
		return
	}

	// Set the upload URL
	uploadUrl := utils.SetUploadUrl(commands.UploadServer)

	//Check that we have an API key set
	if commands.ApiKey == "" {
		log.Error("no API key provided")
		return
	}

	log.Info("uploading to " + uploadUrl)

	switch ctx.Command() {
	case "upload <path>":
		log.Info("Starting Upload")
		//fmt.Println(validation.ValidatePath())
	default:
		println(ctx.Command())
	}

	server.HelloWorld("test")
}