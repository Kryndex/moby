package node

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	apiclient "github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewNodeCommand returns a cobra command for `node` subcommands
func NewNodeCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage Swarm nodes",
		Args:  cli.NoArgs,
		RunE:  dockerCli.ShowHelp,
		Tags:  map[string]string{"version": "1.24"},
	}
	cmd.AddCommand(
		newDemoteCommand(dockerCli),
		newInspectCommand(dockerCli),
		newListCommand(dockerCli),
		newPromoteCommand(dockerCli),
		newRemoveCommand(dockerCli),
		newPsCommand(dockerCli),
		newUpdateCommand(dockerCli),
	)
	return cmd
}

// Reference returns the reference of a node. The special value "self" for a node
// reference is mapped to the current node, hence the node ID is retrieved using
// the `/info` endpoint.
func Reference(ctx context.Context, client apiclient.APIClient, ref string) (string, error) {
	if ref == "self" {
		info, err := client.Info(ctx)
		if err != nil {
			return "", err
		}
		if info.Swarm.NodeID == "" {
			// If there's no node ID in /info, the node probably
			// isn't a manager. Call a swarm-specific endpoint to
			// get a more specific error message.
			_, err = client.NodeList(ctx, types.NodeListOptions{})
			if err != nil {
				return "", err
			}
			return "", errors.New("node ID not found in /info")
		}
		return info.Swarm.NodeID, nil
	}
	return ref, nil
}
