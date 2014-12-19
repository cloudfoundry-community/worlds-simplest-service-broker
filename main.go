package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry-community/types-cf"
	"github.com/go-martini/martini"
	"github.com/kr/pretty"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type serviceBindingResponse struct {
	Credentials    map[string]interface{} `json:"credentials"`
	SyslogDrainURL string                 `json:"syslog_drain_url"`
}

var serviceName, servicePlan, baseGUID, tags string
var serviceBinding serviceBindingResponse

func brokerCatalog() (int, []byte) {
	tagArray := []string{}
	if len(tags) > 0 {
		tagArray = strings.Split(tags, ",")
	}
	catalog := cf.Catalog{
		Services: []*cf.Service{
			{
				ID:          baseGUID + "-service-" + serviceName,
				Name:        serviceName,
				Description: "Shared service for " + serviceName,
				Bindable:    true,
				Tags:        tagArray,
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

func deleteServiceInstance(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)
	return 200, "{}"
}

func createServiceBinding(params martini.Params) (int, []byte) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceName, servicePlan, serviceID)

	json, err := json.Marshal(serviceBinding)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this binding:")
		fmt.Printf("%# v\n", pretty.Formatter(serviceBinding))
		return 500, []byte{}
	}
	return 201, json
}

func deleteServiceBinding(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceName, servicePlan, serviceID)
	return 200, "{}"
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
	credentials := os.Getenv("CREDENTIALS")
	if credentials == "" {
		credentials = "{\"port\": \"4000\"}"
	}
	tags = os.Getenv("TAGS")
	if tags == "" {
		tags = ""
	}
	json.Unmarshal([]byte(credentials), &serviceBinding.Credentials)
	fmt.Printf("%# v\n", pretty.Formatter(serviceBinding))

	m.Get("/v2/catalog", brokerCatalog)
	m.Put("/v2/service_instances/:service_id", createServiceInstance)
	m.Delete("/v2/service_instances/:service_id", deleteServiceInstance)
	m.Put("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)
	m.Delete("/v2/service_instances/:service_id/service_bindings/:binding_id", deleteServiceBinding)

	m.Run()
}
