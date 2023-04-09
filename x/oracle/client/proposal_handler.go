package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/gridiron-zone/gridiron/x/oracle/client/cli"
	"github.com/gridiron-zone/gridiron/x/oracle/client/rest"
)

var (
	RegisterTargetProposalHandler = govclient.NewProposalHandler(cli.NewRegisterTargetProposalCmd, rest.RegisterTargetProposalRESTHandler)
)
