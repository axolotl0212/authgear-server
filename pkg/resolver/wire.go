//+build wireinject

package resolver

import (
	"net/http"

	"github.com/google/wire"

	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/resolver/handler"
	"github.com/authgear/authgear-server/pkg/util/httproute"
)

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	panic(wire.Build(
		deps.RootDependencySet,
		middleware.DependencySet,
		wire.Bind(new(httproute.Middleware), new(*middleware.SentryMiddleware)),
	))
}

func newRootRecoverMiddleware(p *deps.RootProvider) httproute.Middleware {
	panic(wire.Build(
		deps.RootDependencySet,
		middleware.DependencySet,
		wire.Bind(new(httproute.Middleware), new(*middleware.RecoverMiddleware)),
	))
}

func newRequestRecoverMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		DependencySet,
		wire.Bind(new(httproute.Middleware), new(*middleware.RecoverMiddleware)),
	))
}

func newSessionMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		DependencySet,
		wire.Bind(new(httproute.Middleware), new(*session.Middleware)),
	))
}

func newSessionResolveHandler(p *deps.RequestProvider) http.Handler {
	panic(wire.Build(
		DependencySet,
		wire.Bind(new(http.Handler), new(*handler.ResolveHandler)),
	))
}
