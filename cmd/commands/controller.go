package commands

import (
	"github.com/spf13/cobra"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	controllercmd "github.com/nholuongut/argo-events/pkg/reconciler/cmd"
	"github.com/nholuongut/argo-events/pkg/shared/logging"
	sharedutil "github.com/nholuongut/argo-events/pkg/shared/util"
)

func NewControllerCommand() *cobra.Command {
	var (
		leaderElection   bool
		namespaced       bool
		managedNamespace string
		metricsPort      int32
		healthPort       int32
		klogLevel        int
	)

	command := &cobra.Command{
		Use:   "controller",
		Short: "Start the controller",
		Run: func(cmd *cobra.Command, args []string) {
			logging.SetKlogLevel(klogLevel)
			eventOpts := controllercmd.ArgoEventsControllerOpts{
				LeaderElection:   leaderElection,
				ManagedNamespace: managedNamespace,
				Namespaced:       namespaced,
				MetricsPort:      metricsPort,
				HealthPort:       healthPort,
			}
			controllercmd.Start(eventOpts)
		},
	}
	command.Flags().BoolVar(&namespaced, "namespaced", false, "Whether to run in namespaced scope, defaults to false.")
	command.Flags().StringVar(&managedNamespace, "managed-namespace", sharedutil.LookupEnvStringOr("NAMESPACE", "argo-events"), "The namespace that the controller watches when \"--namespaced\" is \"true\".")
	command.Flags().BoolVar(&leaderElection, "leader-election", true, "Enable leader election")
	command.Flags().Int32Var(&metricsPort, "metrics-port", v1alpha1.ControllerMetricsPort, "Metrics port")
	command.Flags().Int32Var(&healthPort, "health-port", v1alpha1.ControllerHealthPort, "Health port")
	command.Flags().IntVar(&klogLevel, "kloglevel", 0, "klog level")
	return command
}
