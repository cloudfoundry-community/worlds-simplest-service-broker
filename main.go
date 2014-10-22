package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-martini/martini"
	"github.com/intel-data/types-cf"
	"github.com/kr/pretty"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

var serviceName, servicePlan, serviceBindingCredentials, baseGUID string

func brokerCatalog() (int, []byte) {
	catalog := cf.Catalog{
		Services: []*cf.Service{
			{
				ID:          baseGUID + "-service-" + serviceName,
				Name:        serviceName,
				Description: "Shared service for " + serviceName,
				Bindable:    true,
				Plans: []*cf.Plan{
					{
						ID:          baseGUID + "-plan-" + servicePlan,
						Name:        servicePlan,
						Description: "Shared service for " + serviceName,
						Free:        true,
					},
				},
			},
		},
	}
	json, err := json.Marshal(catalog)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this catalog:")
		fmt.Printf("%# v\n", pretty.Formatter(catalog))
		return 500, []byte{}
	}
	return 200, json
}

func createServiceInstance(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	fmt.Printf("Creating service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)
	return 201, "{}"
}

func createServiceBinding(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceName, servicePlan, serviceID)
	return 201, serviceBindingCredentials
}

func main() {
	m := martini.Classic()

	baseGUID = os.Getenv("BASE_GUID")
	if baseGUID == "" {
		baseGUID = "29140B3F-0E69-4C7E-8A35"
	}
	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "some-service-name" // replace with cfenv.AppName
	}
	servicePlan = os.Getenv("SERVICE_PLAN")
	if servicePlan == "" {
		servicePlan = "shared"
	}
	serviceBindingCredentials = os.Getenv("CREDENTIALS")
	if serviceBindingCredentials == "" {
		serviceBindingCredentials = "{\"port\": 4000}"
	}

	m.Get("/v2/catalog", brokerCatalog)
	m.Put("/v2/service_instances/:service_id", createServiceInstance)
	m.Put("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)

	m.Run()
}
