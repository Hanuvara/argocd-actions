package main

import (
	"log"
	"os"

	"strconv"

	"github.com/omegion/argocd-actions/internal/argocd"
	ctrl "github.com/omegion/argocd-actions/internal/controller"
)

func main() {
	insecure, err := strconv.ParseBool(os.Getenv("INPUT_INSECURE"))
	if err != nil {
		log.Fatal("error: INPUT_INSECURE failed to parse as bool")
	}

	options := argocd.APIOptions{
		Address:  os.Getenv("INPUT_ADDRESS"),
		Token:    os.Getenv("INPUT_TOKEN"),
		Insecure: insecure,
	}

	api := argocd.NewAPI(options)
	controller := ctrl.NewController(api)

	appName := os.Getenv("INPUT_APPNAME")
	imageTag := os.Getenv("INPUT_IMAGETAG")

	controller.SetImageTag(appName, imageTag)

	err = controller.Sync(appName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
