package main

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/alecthomas/kong"
	"os"
)

func main() {
	var commands struct {
		UploadServer string `help:"Bugsnag On-Premise upload server URL" default:"https://upload.bugsnag.com"`
		Port		 int	`help:"Port number for the upload server" default:"443"`
		ApiKey       string `help:"Bugsnag project API key"`
		Upload       struct {
			// shared options
			Overwrite     bool              `help:"ignore existing upload with same version"`
			Timeout       int               `help:"seconds to wait before failing an upload request"`
			Retries       int               `help:"number of retry attempts before failing a request"`
			UploadOptions map[string]string `help:"additional arguments to pass to the upload request" mapsep:","`

			// required options
			Path []string `help:"Path to directory to search" arg:"" name:"path" type:"path"`
		} `cmd:"" help:"Upload files"`
	}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&commands)

	// Check if we have an apiKey in the request
	if commands.ApiKey == "" {
		log.Error("no API key provided", 1)
	}

	// Check if the path(s) provided are valid.
	if !utils.ValidatePath(commands.Upload.Path) {
		log.Error("path(s) provided is not valid", 1)
	}

	log.Info("uploading files to " + commands.UploadServer)

	// Build a file list form given path(s)
	log.Info("building file list...")

	var fileList []string

	for _, path := range commands.Upload.Path {
		if utils.IsDir(path) {
			log.Info("searching " + " for files...")
			files, err := utils.FilePathWalkDir(path)
			if err != nil {
				log.Error("error getting files from dir", 1)
			}
			for _, s := range files {
				fileList = append(fileList, s)
			}
		} else {
			fileList = append(fileList, path)
		}
	}

	log.Info("File list built..")

	// Build UploadOptions list
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = commands.ApiKey

	for key, value := range commands.Upload.UploadOptions {
		uploadOptions[key] = value
	}

	switch ctx.Command() {

	// Upload command
	case "upload <path>":
		for _, file := range fileList {
			log.Info("starting upload for " + file)
			response, err := upload.All(file, uploadOptions, commands.UploadServer + ":" + string(commands.Port))
			if err != nil {
				log.Error(response, 1)
			}
			log.Info(file + " upload " + response)
		}

	default:
		println(ctx.Command())
	}
}