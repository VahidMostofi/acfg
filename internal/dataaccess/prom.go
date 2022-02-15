package dataaccess

import (
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	log "github.com/sirupsen/logrus"
)

func GetNewPromClientAndQueryAPI(url string) v1.API {
	client, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		log.Error("Error creating client: %v\n", err)
	}

	v1api := v1.NewAPI(client)
	return v1api
}
