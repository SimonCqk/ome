package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ome-agent",
	Short: "Run OME Agent",
	Long:  "OME Agent is a swiss army knife for OME inference service, training job, model management, etc.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Register all agent commands
	rootCmd.AddCommand(CreateAgentCommand(NewEnigmaAgent()))
	rootCmd.AddCommand(CreateAgentCommand(NewHFDownloadAgent()))
	rootCmd.AddCommand(CreateAgentCommand(NewReplicaAgent()))
	rootCmd.AddCommand(CreateAgentCommand(NewServingAgent()))
	rootCmd.AddCommand(CreateAgentCommand(NewFineTunedAdapterAgent()))
}
