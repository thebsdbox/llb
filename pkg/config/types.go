package config

import "net/url"

// Services contains all defined services
var Services []Service

// Service defines the external connectivity and the backend connections
type Service struct {
	Name        string     `json:"name"`       // The Name for the Service
	Frontend    FrontEnd   `json:"frontend"`   // A defined front end
	Endpoint    []Endpoint `json:"endpoints"`  // The multiple endpoints
	ServiceType string     `json:"type"`       // Type of service being exposed (should match the endpoint type)
	Persistent  bool       `json:"persistent"` // Keep persistent states
}

// FrontEnd is the exposed load balancer/service
type FrontEnd struct {
	Adapter     string // Adapter to bind to (optional)
	Name        string // Name of the frontend
	Port        int    // Exposed service
	Address     string // Adress binding
	ServiceType string // Type of service being exposed (should match the endpoint type)
	parsedurl   *url.URL
}

// Endpoint is a single endpoint that will be load balanced over
type Endpoint struct {
	Name        string
	ServiceType string
	RawEndpoint string
	Address     string
	Port        int
	Connections int
	parsedurl   *url.URL
}
