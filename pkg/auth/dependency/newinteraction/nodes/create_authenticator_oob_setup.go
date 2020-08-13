package nodes

import (
	"github.com/authgear/authgear-server/pkg/auth/dependency/authenticator"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity"
	"github.com/authgear/authgear-server/pkg/auth/dependency/newinteraction"
	"github.com/authgear/authgear-server/pkg/lib/authn"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/otp"
)

func init() {
	newinteraction.RegisterNode(&NodeCreateAuthenticatorOOBSetup{})
}

type InputCreateAuthenticatorOOBSetup interface {
	GetOOBChannel() authn.AuthenticatorOOBChannel
	GetOOBTarget() string
}

type EdgeCreateAuthenticatorOOBSetup struct {
	Stage newinteraction.AuthenticationStage
	Tag   []string

	// Either have Channel and Target
	Channel authn.AuthenticatorOOBChannel
	Target  string
	// Or have AllowedChannels
	AllowedChannels []authn.AuthenticatorOOBChannel
}

func (e *EdgeCreateAuthenticatorOOBSetup) AuthenticatorType() authn.AuthenticatorType {
	return authn.AuthenticatorTypeOOB
}

func (e *EdgeCreateAuthenticatorOOBSetup) HasDefaultTag() bool {
	return false
}

func (e *EdgeCreateAuthenticatorOOBSetup) Instantiate(ctx *newinteraction.Context, graph *newinteraction.Graph, rawInput interface{}) (newinteraction.Node, error) {
	var target string
	var channel authn.AuthenticatorOOBChannel

	if e.Channel != "" && e.Target != "" {
		channel = e.Channel
		target = e.Target
	} else {
		input, ok := rawInput.(InputCreateAuthenticatorOOBSetup)
		if !ok {
			return nil, newinteraction.ErrIncompatibleInput
		}
		channel = input.GetOOBChannel()
		if channel == "" {
			return nil, newinteraction.ErrIncompatibleInput
		}
		target = input.GetOOBTarget()
	}

	var spec *authenticator.Spec
	var identityInfo *identity.Info
	if e.Stage == newinteraction.AuthenticationStagePrimary {
		// Primary OOB authenticators must be bound to login ID identity
		identityInfo = graph.MustGetUserLastIdentity()
		if identityInfo.Type != authn.IdentityTypeLoginID {
			panic("interaction: OOB authenticator identity must be login ID")
		}

		spec = &authenticator.Spec{
			UserID: identityInfo.UserID,
			Tag:    stageToAuthenticatorTag(e.Stage),
			Type:   authn.AuthenticatorTypeOOB,
			Props:  map[string]interface{}{},
		}

		// Ignore given OOB target, use channel & target inferred from identity
		loginIDKey := identityInfo.Claims[identity.IdentityClaimLoginIDKey].(string)
		for _, t := range ctx.Config.Identity.LoginID.Keys {
			if t.Key != loginIDKey {
				continue
			}
			switch t.Type {
			case config.LoginIDKeyTypeEmail:
				channel = authn.AuthenticatorOOBChannelEmail
			case config.LoginIDKeyTypePhone:
				channel = authn.AuthenticatorOOBChannelSMS
			default:
				return nil, newinteraction.ConfigurationViolated.New("this login ID cannot be used for OOB authentication")
			}
			break
		}
		target = identityInfo.Claims[identity.IdentityClaimLoginIDValue].(string)

	} else {
		userID := graph.MustGetUserID()
		spec = &authenticator.Spec{
			UserID: userID,
			Tag:    stageToAuthenticatorTag(e.Stage),
			Type:   authn.AuthenticatorTypeOOB,
			Props:  map[string]interface{}{},
		}

		// Normalize the target.
		switch channel {
		case authn.AuthenticatorOOBChannelEmail:
			var err error
			target, err = ctx.LoginIDNormalizerFactory.NormalizerWithLoginIDType(config.LoginIDKeyTypeEmail).Normalize(target)
			if err != nil {
				return nil, err
			}
		case authn.AuthenticatorOOBChannelSMS:
			var err error
			target, err = ctx.LoginIDNormalizerFactory.NormalizerWithLoginIDType(config.LoginIDKeyTypePhone).Normalize(target)
			if err != nil {
				return nil, err
			}
		}
	}

	spec.Tag = append(spec.Tag, e.Tag...)

	spec.Props[authenticator.AuthenticatorPropOOBOTPChannelType] = string(channel)
	switch channel {
	case authn.AuthenticatorOOBChannelSMS:
		spec.Props[authenticator.AuthenticatorPropOOBOTPPhone] = target
	case authn.AuthenticatorOOBChannelEmail:
		spec.Props[authenticator.AuthenticatorPropOOBOTPEmail] = target
	}

	info, err := ctx.Authenticators.New(spec, "")
	if err != nil {
		return nil, err
	}

	secret, err := otp.GenerateTOTPSecret()
	if err != nil {
		return nil, err
	}

	result, err := sendOOBCode(ctx, e.Stage, false, info, secret)
	if err != nil {
		return nil, err
	}

	return &NodeCreateAuthenticatorOOBSetup{
		Stage:           e.Stage,
		AllowedChannels: e.AllowedChannels,
		Authenticator:   info,
		Secret:          secret,
		Target:          target,
		Channel:         result.Channel,
		CodeLength:      result.CodeLength,
		SendCooldown:    result.SendCooldown,
	}, nil
}

type NodeCreateAuthenticatorOOBSetup struct {
	Stage           newinteraction.AuthenticationStage `json:"stage"`
	AllowedChannels []authn.AuthenticatorOOBChannel    `json:"allowed_channels"`
	Authenticator   *authenticator.Info                `json:"authenticator"`
	Secret          string                             `json:"secret"`
	Target          string                             `json:"target"`
	Channel         string                             `json:"channel"`
	CodeLength      int                                `json:"code_length"`
	SendCooldown    int                                `json:"send_cooldown"`
}

// GetOOBOTPTarget implements OOBOTPNode.
func (n *NodeCreateAuthenticatorOOBSetup) GetOOBOTPTarget() string {
	return n.Target
}

// GetOOBOTPChannel implements OOBOTPNode.
func (n *NodeCreateAuthenticatorOOBSetup) GetOOBOTPChannel() string {
	return n.Channel
}

// GetOOBOTPCodeSendCooldown implements OOBOTPNode.
func (n *NodeCreateAuthenticatorOOBSetup) GetOOBOTPCodeSendCooldown() int {
	return n.SendCooldown
}

// GetOOBOTPCodeLength implements OOBOTPNode.
func (n *NodeCreateAuthenticatorOOBSetup) GetOOBOTPCodeLength() int {
	return n.CodeLength
}

func (n *NodeCreateAuthenticatorOOBSetup) Prepare(ctx *newinteraction.Context, graph *newinteraction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorOOBSetup) Apply(perform func(eff newinteraction.Effect) error, graph *newinteraction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorOOBSetup) DeriveEdges(graph *newinteraction.Graph) ([]newinteraction.Edge, error) {
	return []newinteraction.Edge{
		&EdgeOOBResendCode{
			Stage:            n.Stage,
			IsAuthenticating: false,
			Authenticator:    n.Authenticator,
			Secret:           n.Secret,
		},
		&EdgeCreateAuthenticatorOOB{Stage: n.Stage, Authenticator: n.Authenticator, Secret: n.Secret},
	}, nil
}
