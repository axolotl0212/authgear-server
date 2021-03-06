package nodes

import (
	"errors"

	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator"
	"github.com/authgear/authgear-server/pkg/lib/interaction"
)

func init() {
	interaction.RegisterNode(&NodeCreateAuthenticatorTOTP{})
}

type InputCreateAuthenticatorTOTP interface {
	GetTOTP() string
	GetTOTPDisplayName() string
}

type EdgeCreateAuthenticatorTOTP struct {
	Stage         interaction.AuthenticationStage
	Authenticator *authenticator.Info
}

func (e *EdgeCreateAuthenticatorTOTP) Instantiate(ctx *interaction.Context, graph *interaction.Graph, rawInput interface{}) (interaction.Node, error) {
	input, ok := rawInput.(InputCreateAuthenticatorTOTP)
	if !ok {
		return nil, interaction.ErrIncompatibleInput
	}

	info := cloneAuthenticator(e.Authenticator)
	info.Claims[authenticator.AuthenticatorClaimTOTPDisplayName] = input.GetTOTPDisplayName()

	err := ctx.Authenticators.VerifySecret(info, nil, input.GetTOTP())
	if errors.Is(err, authenticator.ErrInvalidCredentials) {
		return nil, interaction.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return &NodeCreateAuthenticatorTOTP{Stage: e.Stage, Authenticator: info}, nil
}

type NodeCreateAuthenticatorTOTP struct {
	Stage         interaction.AuthenticationStage `json:"stage"`
	Authenticator *authenticator.Info             `json:"authenticator"`
}

func (n *NodeCreateAuthenticatorTOTP) Prepare(ctx *interaction.Context, graph *interaction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorTOTP) Apply(perform func(eff interaction.Effect) error, graph *interaction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorTOTP) DeriveEdges(graph *interaction.Graph) ([]interaction.Edge, error) {
	return []interaction.Edge{
		&EdgeCreateAuthenticatorEnd{
			Stage:          n.Stage,
			Authenticators: []*authenticator.Info{n.Authenticator},
		},
	}, nil
}
