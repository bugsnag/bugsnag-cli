package utils

const DefaultUploadURL = "https://upload.bugsnag.com/"

func SetUploadUrl(url string) string {
	if url == "" {
		url = DefaultUploadURL
	}
	return url
}
