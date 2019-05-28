package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var endPointIndex int // Holds the previous endpoint (for determining decisions on next endpoint)

func init() {
	// Start the index negative as it will be incrememnted of first approach
	endPointIndex = -1
}

// ValidateEndpointURLS will run through the endpoints and ensure that they're a valid URL
func ValidateEndpointURLS(endpoints *[]Endpoint) error {
	// Find the service

	for i := range *endpoints {
		log.Debugf("Parsing [%s]", (*endpoints)[i].RawEndpoint)
		u, err := url.Parse((*endpoints)[i].RawEndpoint)
		if err != nil {
			return err
		}

		// No error is returned if the prefix/schema is missing
		// If the Host is empty then we were unable to parse correctly (could be prefix is missing)
		if u.Host == "" {
			return fmt.Errorf("Unable to parse [%s], ensure it's prefixed with http(s)://", (*endpoints)[i].RawEndpoint)
		}
		(*endpoints)[i].Address = u.Hostname()
		// if a port is specified then update the internal endpoint stuct, if not rely on the schema
		if u.Port() != "" {
			portNum, err := strconv.Atoi(u.Port())
			if err != nil {
				return err
			}
			(*endpoints)[i].Port = portNum
		}
		(*endpoints)[i].parsedurl = u
	}
	return nil
}

// ValidateEndpoints will run through the endpoints and ensure that they're a valid endpoint
func ValidateEndpoints(endpoints *[]Endpoint) error {
	for i := range *endpoints {
		log.Debugf("Parsing [%s]", (*endpoints)[i].RawEndpoint)
		host, port, err := net.SplitHostPort((*endpoints)[i].RawEndpoint)
		if err != nil {
			return err
		}
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		(*endpoints)[i].Address = host
		(*endpoints)[i].Port = portNum
	}
	return nil
}

// ReturnEndpointURL - returns an endpoint
func (s Service) ReturnEndpointURL() *url.URL {

	if endPointIndex != len(s.Endpoint)-1 {
		endPointIndex++
	} else {
		// reset the index to the beginning
		endPointIndex = 0
	}
	// TODO - weighting, decision algorythmn
	return s.Endpoint[endPointIndex].parsedurl
}

// ReturnEndpointAddr - returns an endpoint
func (s Service) ReturnEndpointAddr() string {

	if endPointIndex != len(s.Endpoint)-1 {
		endPointIndex++
	} else {
		// reset the index to the beginning
		endPointIndex = 0
	}
	// TODO - weighting, decision algorythmn
	return fmt.Sprintf("%s:%d", s.Endpoint[endPointIndex].Address, s.Endpoint[endPointIndex].Port)
}
