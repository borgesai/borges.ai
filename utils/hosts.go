package utils

import (
	"net/http"
	"os"
)

func IsCustomHost(req *http.Request) bool{
	devHost := "localhost:" + os.Getenv("PORT")
	return (req.Host != os.Getenv("SERVICE_HOST")) && (req.Host != devHost)
}
