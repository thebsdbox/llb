package config

// Configuration for exposing a load balancer service

// Address -
func (s Service) Address() string {
	return s.Frontend.Address
}

// Port --
func (s Service) Port() int {
	return s.Frontend.Port
}
