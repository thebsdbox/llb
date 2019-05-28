package servertcp

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/thebsdbox/llb/pkg/config"
)

// 1. Load balancer port is exposed
// 2. We listen
// 3. On connection we connect to an endpoint
// [loop]
// 4. We read from the load balancer port
// 5. We write traffic to the endpoint
// 6. We read response from endpoint
// 7. We write response to load balancer
// [goto loop]

func persistentConnection(frontendConnection net.Conn, s *config.Service) error {
	// Connect to Endpoint
	ep := s.ReturnEndpointAddr()

	log.Debugf("Attempting endpoint [%s]", ep)

	endpoint, err := net.Dial("tcp", ep)
	if err != nil {
		fmt.Println("dial error:", err)
		// return nil, err
	}
	log.Debugf("succesfully connected to [%s]", ep)

	// Build endpoint <front end> connectivity
	go func() { io.Copy(frontendConnection, endpoint) }()
	go func() { io.Copy(endpoint, frontendConnection) }()

	return nil
}
