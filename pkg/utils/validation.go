package utils

import "os"

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
