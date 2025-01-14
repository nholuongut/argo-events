package commands

import (
	"github.com/spf13/cobra"

	sensorcmd "github.com/nholuongut/argo-events/pkg/sensors/cmd"
)

func NewSensorCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "sensor-service",
		Short: "Start a Sensor service",
		Run: func(cmd *cobra.Command, args []string) {
			sensorcmd.Start()
		},
	}
	return command
}
