package bucket

import "os"

func getCredentials() string {
	return os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
}

func getBucketName() string {
	return os.Getenv("GOOGLE_BUCKET_NAME")
}
