package main

import (
	"log"
	"os"

	"strconv"

	"github.com/omegion/argocd-actions/internal/argocd"
	ctrl "github.com/omegion/argocd-actions/internal/controller"
)

func main() {
	insecure, err := strconv.ParseBool(os.Getenv("IS_INSECURE"))
	if err != nil {
		log.Fatal("error: IS_INSECURE failed to parse as bool")
	}

	options := argocd.APIOptions{
		Address:  os.Getenv("INPUT_ADDRESS"),
		Token:    os.Getenv("INPUT_TOKEN"),
		Insecure: insecure,
	}

	api := argocd.NewAPI(options)
	controller := ctrl.NewController(api)

	err = controller.Sync(os.Getenv("INPUT_APPNAME"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
