package utils

func SetUploadUrl(url string) string {
	if url == "" {
		url = "https://upload.bugsnag.com/"
	}
	return url
}
