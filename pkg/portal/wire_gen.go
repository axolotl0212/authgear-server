// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package portal

import (
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/portal/deps"
	"github.com/authgear/authgear-server/pkg/portal/graphql"
	"github.com/authgear/authgear-server/pkg/portal/loader"
	"github.com/authgear/authgear-server/pkg/portal/service"
	"github.com/authgear/authgear-server/pkg/portal/session"
	"github.com/authgear/authgear-server/pkg/portal/transport"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"net/http"
)

// Injectors from wire.go:

func newRecoverMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	factory := rootProvider.LoggerFactory
	recoveryLogger := middleware.NewRecoveryLogger(factory)
	recoverMiddleware := &middleware.RecoverMiddleware{
		Logger: recoveryLogger,
	}
	return recoverMiddleware
}

func newSentryMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	hub := rootProvider.SentryHub
	environmentConfig := rootProvider.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newSessionInfoMiddleware(p *deps.RequestProvider) httproute.Middleware {
	sessionMiddleware := &session.Middleware{}
	return sessionMiddleware
}

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	request := p.Request
	context := deps.ProvideRequestContext(request)
	viewerLoader := &loader.ViewerLoader{
		Context: context,
	}
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	configGetter := &deps.ConfigGetter{
		Request:      request,
		ConfigSource: configSource,
	}
	appService := &service.AppService{
		ConfigGetter: configGetter,
	}
	appLoader := &loader.AppLoader{
		Apps: appService,
	}
	graphqlContext := &graphql.Context{
		Viewer: viewerLoader,
		Apps:   appLoader,
	}
	graphQLHandler := &transport.GraphQLHandler{
		DevMode:        devMode,
		GraphQLContext: graphqlContext,
	}
	return graphQLHandler
}

func newRuntimeConfigHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	authgearConfig := rootProvider.AuthgearConfig
	runtimeConfigHandler := &transport.RuntimeConfigHandler{
		AuthgearConfig: authgearConfig,
	}
	return runtimeConfigHandler
}

func newAdminAPIHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	adminAPIConfig := rootProvider.AdminAPIConfig
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	clock := _wireSystemClockValue
	adder := &authz.Adder{
		Clock: clock,
	}
	adminAPIService := &service.AdminAPIService{
		AdminAPIConfig: adminAPIConfig,
		ConfigSource:   configSource,
		AuthzAdder:     adder,
	}
	adminAPIHandler := &transport.AdminAPIHandler{
		ConfigResolver:   adminAPIService,
		EndpointResolver: adminAPIService,
		AuthzAdder:       adminAPIService,
	}
	return adminAPIHandler
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)