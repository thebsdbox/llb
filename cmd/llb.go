package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/llb/pkg/config"
	"github.com/thebsdbox/llb/pkg/nft"
	"github.com/thebsdbox/llb/pkg/serverhttp"
	"github.com/thebsdbox/llb/pkg/servertcp"
)

// Variable start below

// Release - this struct contains the release information populated when building plunder
var Release struct {
	Version string
	Build   string
}

// This manages the level of output
var logLevel int

// service name , if created from the CLI
var serviceName string

// service type , if created from the CLI
var serviceType string

// exposed port, if created from CLI
var exposedPort int

// endpoints contains all of the Load balancer end points
var endpoints []string

var llbCmd = &cobra.Command{
	Use:   "llb",
	Short: "The Little Load Balancer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var llbNftCmd = &cobra.Command{
	Use:   "nft",
	Short: "The Little Load Balancer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		nft.CreateLoadBalancer()
	},
}

var llbServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a Load Balancer Server",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		// Validate the CLI flags

		// Endpoints are defined from CLI
		if len(endpoints) != 0 {
			err := config.ParseCLI(serviceName, serviceType, endpoints, exposedPort)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}

		if config.ConfigPath != "" {
			err := config.ParseConfig()
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
		serverhttp.Start()
		servertcp.Start()
	},
}

var llbVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version and Release Info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("LB Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}

var llbExample = &cobra.Command{
	Use:   "example",
	Short: "Print Example configuration",
	Run: func(cmd *cobra.Command, args []string) {
		exampleConfig := config.Service{
			Name:        "Example Service 1",
			ServiceType: "tcp",
			Frontend: config.FrontEnd{
				Adapter: "eth0",
				Port:    6443,
				Address: "127.0.0.1",
			},
		}

		//Create endpoint 1
		exampleEndpoint1 := config.Endpoint{
			Name:    "K8s-master01",
			Address: "192.168.0.100",
			Port:    6443,
		}
		exampleConfig.Endpoint = append(exampleConfig.Endpoint, exampleEndpoint1)

		// Create endpoint 2
		exampleEndpoint2 := config.Endpoint{
			Name:    "K8s-master02",
			Address: "192.168.0.101",
			Port:    6443,
		}
		exampleConfig.Endpoint = append(exampleConfig.Endpoint, exampleEndpoint2)

		config.Services = append(config.Services, exampleConfig)

		b, _ := yaml.Marshal(config.Services)
		fmt.Printf("%s\n", b)
		return
	},
}

// Functions start below

func init() {
	// Parse flags
	llbServerCmd.Flags().StringSliceVarP(&endpoints, "endpoint", "e", []string{}, "The definition of one (or more) endpoints")
	llbServerCmd.Flags().StringVarP(&config.ConfigPath, "config", "c", "", "Path to a yaml configuration file")
	llbServerCmd.Flags().StringVarP(&serviceName, "service", "s", "default", "The Name of a service to create")
	llbServerCmd.Flags().StringVarP(&serviceType, "type", "t", "", "The Type of service being exposed [HTTP/TCP/UDP]")

	// Local port that this service is exposed on
	llbServerCmd.Flags().IntVarP(&exposedPort, "port", "p", 0, "The port this service will be exposed on")

	llbCmd.PersistentFlags().IntVarP(&logLevel, "loglevel", "l", 4, "Level of logging messages")

	// Add subcommands
	llbCmd.AddCommand(llbExample)
	llbCmd.AddCommand(llbNftCmd)
	llbCmd.AddCommand(llbServerCmd)
	llbCmd.AddCommand(llbVersionCmd)
}

// Execute - starts the command parsing process
func Execute() {
	// TODO - Logging from environment variable is broken
	if os.Getenv("LB_LOGLVL") != "" {
		i, err := strconv.ParseInt(os.Getenv("LB_LOGLVL"), 10, 8)
		if err != nil {
			log.Fatalf("Error parsing environment variable [LB_LOGLVL")
		}
		// We've only parsed to an 8bit integer, however i is still a int64 so needs casting
		logLevel = int(i)
	} else {
		// Default to logging anything Info and below
		logLevel = int(log.InfoLevel)
	}

	if err := llbCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
