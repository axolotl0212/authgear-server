// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package admin

import (
	"github.com/authgear/authgear-server/pkg/admin/graphql"
	"github.com/authgear/authgear-server/pkg/admin/loader"
	"github.com/authgear/authgear-server/pkg/admin/transport"
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/oob"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/password"
	service2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/totp"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/anonymous"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/loginid"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"net/http"
)

// Injectors from wire.go:

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	hub := p.SentryHub
	environmentConfig := p.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newRootRecoverMiddleware(p *deps.RootProvider) httproute.Middleware {
	factory := p.LoggerFactory
	recoveryLogger := middleware.NewRecoveryLogger(factory)
	recoverMiddleware := &middleware.RecoverMiddleware{
		Logger: recoveryLogger,
	}
	return recoverMiddleware
}

func newRequestRecoverMiddleware(p *deps.RequestProvider) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	recoveryLogger := middleware.NewRecoveryLogger(factory)
	recoverMiddleware := &middleware.RecoverMiddleware{
		Logger: recoveryLogger,
	}
	return recoverMiddleware
}

func newAuthorizationMiddleware(p *deps.RequestProvider, auth config.AdminAPIAuth) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	logger := authz.NewLogger(factory)
	configConfig := appProvider.Config
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	secretConfig := configConfig.SecretConfig
	adminAPIAuthKey := deps.ProvideAdminAPIAuthKeyMaterials(secretConfig)
	clock := _wireSystemClockValue
	authzMiddleware := &authz.Middleware{
		Logger:  logger,
		Auth:    auth,
		AppID:   appID,
		AuthKey: adminAPIAuthKey,
		Clock:   clock,
	}
	return authzMiddleware
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	appProvider := p.AppProvider
	configConfig := appProvider.Config
	secretConfig := configConfig.SecretConfig
	databaseCredentials := deps.ProvideDatabaseCredentials(secretConfig)
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	sqlBuilder := db.ProvideSQLBuilder(databaseCredentials, appID)
	request := p.Request
	context := deps.ProvideRequestContext(request)
	handle := appProvider.Database
	sqlExecutor := db.SQLExecutor{
		Context:  context,
		Database: handle,
	}
	store := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticationConfig := appConfig.Authentication
	identityConfig := appConfig.Identity
	serviceStore := &service.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	loginidStore := &loginid.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	loginIDConfig := identityConfig.LoginID
	rootProvider := appProvider.RootProvider
	reservedNameChecker := rootProvider.ReservedNameChecker
	typeCheckerFactory := &loginid.TypeCheckerFactory{
		Config:              loginIDConfig,
		ReservedNameChecker: reservedNameChecker,
	}
	checker := &loginid.Checker{
		Config:             loginIDConfig,
		TypeCheckerFactory: typeCheckerFactory,
	}
	normalizerFactory := &loginid.NormalizerFactory{
		Config: loginIDConfig,
	}
	clockClock := _wireSystemClockValue
	provider := &loginid.Provider{
		Store:             loginidStore,
		Config:            loginIDConfig,
		Checker:           checker,
		NormalizerFactory: normalizerFactory,
		Clock:             clockClock,
	}
	oauthStore := &oauth.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	oauthProvider := &oauth.Provider{
		Store: oauthStore,
		Clock: clockClock,
	}
	anonymousStore := &anonymous.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	anonymousProvider := &anonymous.Provider{
		Store: anonymousStore,
		Clock: clockClock,
	}
	serviceService := &service.Service{
		Authentication: authenticationConfig,
		Identity:       identityConfig,
		Store:          serviceStore,
		LoginID:        provider,
		OAuth:          oauthProvider,
		Anonymous:      anonymousProvider,
	}
	factory := appProvider.LoggerFactory
	logger := verification.NewLogger(factory)
	verificationConfig := appConfig.Verification
	store2 := &service2.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	passwordStore := &password.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorConfig := appConfig.Authenticator
	authenticatorPasswordConfig := authenticatorConfig.Password
	passwordLogger := password.NewLogger(factory)
	historyStore := &password.HistoryStore{
		Clock:       clockClock,
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	passwordChecker := password.ProvideChecker(authenticatorPasswordConfig, historyStore)
	queue := appProvider.TaskQueue
	passwordProvider := &password.Provider{
		Store:           passwordStore,
		Config:          authenticatorPasswordConfig,
		Clock:           clockClock,
		Logger:          passwordLogger,
		PasswordHistory: historyStore,
		PasswordChecker: passwordChecker,
		TaskQueue:       queue,
	}
	totpStore := &totp.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorTOTPConfig := authenticatorConfig.TOTP
	totpProvider := &totp.Provider{
		Store:  totpStore,
		Config: authenticatorTOTPConfig,
		Clock:  clockClock,
	}
	authenticatorOOBConfig := authenticatorConfig.OOB
	oobStore := &oob.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	oobProvider := &oob.Provider{
		Config: authenticatorOOBConfig,
		Store:  oobStore,
		Clock:  clockClock,
	}
	service3 := &service2.Service{
		Store:    store2,
		Password: passwordProvider,
		TOTP:     totpProvider,
		OOBOTP:   oobProvider,
	}
	redisHandle := appProvider.Redis
	storeRedis := &verification.StoreRedis{
		Redis: redisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	verificationService := &verification.Service{
		Logger:         logger,
		Config:         verificationConfig,
		Clock:          clockClock,
		Identities:     serviceService,
		Authenticators: service3,
		Store:          storeRedis,
	}
	queries := &user.Queries{
		Store:        store,
		Identities:   serviceService,
		Verification: verificationService,
	}
	userLoader := &loader.UserLoader{
		Users: queries,
	}
	identityLoader := &loader.IdentityLoader{
		Identities: serviceService,
	}
	authenticatorLoader := &loader.AuthenticatorLoader{
		Authenticators: service3,
	}
	graphqlContext := &graphql.Context{
		Users:          userLoader,
		Identities:     identityLoader,
		Authenticators: authenticatorLoader,
	}
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	graphQLHandler := &transport.GraphQLHandler{
		GraphQLContext: graphqlContext,
		DevMode:        devMode,
		Database:       handle,
	}
	return graphQLHandler
}
