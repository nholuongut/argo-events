package commands

import (
	"github.com/spf13/cobra"

	eventsourcecmd "github.com/nholuongut/argo-events/pkg/eventsources/cmd"
)

func NewEventSourceCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "eventsource-service",
		Short: "Start an EventSource service",
		Run: func(cmd *cobra.Command, args []string) {
			eventsourcecmd.Start()
		},
	}
	return command
}
