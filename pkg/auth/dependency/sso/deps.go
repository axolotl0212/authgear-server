package sso

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(StateCodec), "*"),
	wire.Struct(new(UserInfoDecoder), "*"),
	wire.Struct(new(OAuthProviderFactory), "*"),
)
