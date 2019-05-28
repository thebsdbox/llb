package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

// ConfigPath points to a configuration file
var ConfigPath string

// ParseConfig will attempt to read through the configuration file and ensure it's valid
func ParseConfig() error {
	if ConfigPath != "" {
		log.Infof("Reading deployment configuration from [%s]", ConfigPath)

		// // Check the actual path from the string
		if _, err := os.Stat(ConfigPath); !os.IsNotExist(err) {
			b, err := ioutil.ReadFile(ConfigPath)
			if err != nil {
				log.Fatalf("%v", err)
			}

			jsonBytes, err := yaml.YAMLToJSON(b)
			if err == nil {
				// If there were no errors then the YAML => JSON was succesful, no attempt to unmarshall
				err = json.Unmarshal(jsonBytes, &Services)
				if err != nil {
					return fmt.Errorf("%v", err)
				}

			} else {
				// Couldn't parse the yaml to JSON
				// Attempt to parse it as JSON
				err = json.Unmarshal(b, &Services)
				if err != nil {
					return fmt.Errorf("%v", err)
				}
			}
		}
	}

	return nil
}

// ParseCLI will attempt to read through the CLI and ensure it's valid
func ParseCLI(serviceName, serviceType string, endpoints []string, port int) error {
	log.Debugf("Parsing CLI configuration for service named [%s]", serviceName)

	// Build endpoint array
	var e []Endpoint
	for i := range endpoints {
		newEndpoint := Endpoint{
			Name:        fmt.Sprintf("%s-%d", serviceName, i),
			ServiceType: serviceType,
			RawEndpoint: endpoints[i],
		}
		e = append(e, newEndpoint)
	}
	switch serviceType {
	case strings.ToLower("http"):
		err := ValidateEndpointURLS(&e)
		if err != nil {
			return err
		}
		// Display configuration
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "Name\tHost\tPort\n")
		for y := range e {
			fmt.Fprintf(w, "%s\t%s\t%d\n", e[y].Name, e[y].Address, e[y].Port)
		}
		w.Flush()

		// Build a new service
		s := Service{
			Name:        serviceName,
			Endpoint:    e,
			ServiceType: serviceType,
			Frontend: FrontEnd{
				Name:        fmt.Sprintf("%s-fe", serviceName),
				Port:        port,
				ServiceType: serviceType,
			},
		}

		// Add load balancing service to the services array
		Services = append(Services, s)

	case strings.ToLower("udp"):
		// TODO (dan)
	case strings.ToLower("tcp"):
		err := ValidateEndpoints(&e)
		if err != nil {
			return err
		}

		// Display configuration
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "Name\tHost\tPort\n")
		for y := range e {
			fmt.Fprintf(w, "%s\t%s\t%d\n", e[y].Name, e[y].Address, e[y].Port)
		}
		w.Flush()

		s := Service{
			Name:        serviceName,
			Endpoint:    e,
			ServiceType: serviceType,
			Frontend: FrontEnd{
				Name:        fmt.Sprintf("%s-fe", serviceName),
				Port:        port,
				ServiceType: serviceType,
			},
		}
		// Add load balancing service to the services array
		Services = append(Services, s)
	default:
		return fmt.Errorf("Unknow Service type [%s]", serviceType)
	}
	return nil
}
