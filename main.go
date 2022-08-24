package main

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/alecthomas/kong"
	"os"
	"strconv"
)

type UploadPaths []string

// Validate that the path(s) exist
func (p UploadPaths) Validate() error {
	for _,path := range p {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var commands struct {
		UploadAPIRootUrl string `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
		Port		 int	`help:"Port number for the upload server" default:"443"`
		ApiKey       string `help:"Bugsnag project API key"`
		Upload       struct {
			// shared options
			Overwrite     bool              `help:"ignore existing upload with same version"`
			Timeout       int               `help:"seconds to wait before failing an upload request" default:"300"`
			Retries       int               `help:"number of retry attempts before failing a request" default:"0"`
			UploadOptions map[string]string `help:"additional arguments to pass to the upload request" mapsep:","`

			// required options
			Path UploadPath `help:"Path to directory to search" arg:"" name:"path" type:"path"`
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

	// Build connection URI
	endpoint := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	log.Info("uploading files to " + endpoint)

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

	if commands.Upload.Overwrite {
		uploadOptions["overwrite"] =  "true"
	}

	uploadOptions["retries"] =  strconv.Itoa(commands.Upload.Retries)

	for key, value := range commands.Upload.UploadOptions {
		uploadOptions[key] = value
	}

	switch ctx.Command() {

	// Upload command
	case "upload <path>":
		for _, file := range fileList {
			log.Info("starting upload for " + file)
			response, err := upload.All(file, uploadOptions, endpoint, commands.Upload.Timeout)
			if err != nil {
				fmt.Println(err)
				log.Error(response, 1)
			}
			log.Success(file + " upload " + response)
		}

	default:
		println(ctx.Command())
	}
}