package main

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/alecthomas/kong"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func main() {
	defer log.Stop()

	var commands struct {
		UploadServer string `help:"Bugsnag On-Premise upload server URL"`
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

	ctx := kong.Parse(&commands)

	// Check if we have an apiKey in the request
	if commands.ApiKey == "" {
		log.Error("no API key provided")
		utils.CleanupAndExit(1)
	}

	// Check if the path(s) provided are valid.
	if !utils.ValidatePath(commands.Upload.Path) {
		log.Error("path(s) provided is not valid")
		utils.CleanupAndExit(1)
	}

	// Set the server upload URL
	uploadUrl := utils.SetUploadUrl(commands.UploadServer)

	log.Info("uploading files to " + uploadUrl)

	// Build a file list form given path(s)
	log.Info("building file list...")

	var fileList []string

	for _, path := range commands.Upload.Path {
		if utils.IsDir(path) {
			log.Info("searching " + " for files...")
			files, err := utils.FilePathWalkDir(path)
			if err != nil {
				log.Error("error getting files from dir")
				utils.CleanupAndExit(1)
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

	fmt.Println(ctx.Command())

	switch ctx.Command() {

	// Upload command
	case "upload <path>":
		for _, file := range fileList {
			log.Info("starting upload for " + file)
			response, err := upload.All(file, uploadOptions, uploadUrl)
			if err != nil {
				log.Error(response)
				utils.CleanupAndExit(1)
			}
			log.Info(file + " upload " + response)
		}

	default:
		println(ctx.Command())
	}
}