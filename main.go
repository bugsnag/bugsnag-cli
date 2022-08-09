package main

import (
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/validation"
	"github.com/alecthomas/kong"
	"os"
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
			Overwrite bool   	`help:"ignore existing upload with same version"`
			Timeout   int    	`help:"seconds to wait before failing an upload request"`
			Retries   int    	`help:"number of retry attempts before failing a request"`

			// required options
			Path 	  []string 	`arg:"" name:"path" help:"Path to directory to search" type:"path"`
		} `cmd:"" help:"Upload files"`
	}

	ctx := kong.Parse(&commands)

	// Check if path(s) provided are
	for _,path := range commands.Upload.Path {
		log.Info("Checking if %s is a valid file or path", path)
		if !validation.ValidatePath(path) {
			log.Error("%s is not valid file or path", path)
			os.Exit(5)
		}
	}

	switch ctx.Command() {
	case "upload <path>":
		log.Info("Starting Upload")
		//fmt.Println(validation.ValidatePath())
	default:
		println(ctx.Command())
	}

	server.HelloWorld("test")
}