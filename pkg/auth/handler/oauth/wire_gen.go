// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package oauth

import (
	"github.com/skygeario/skygear-server/pkg/auth"
	auth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth"
	redis3 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/bearertoken"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/oob"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/password"
	provider2 "github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/recoverycode"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/totp"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/challenge"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/anonymous"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/flows"
	redis2 "github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/redis"
	oauth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oidc"
	handler2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oidc/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	redis4 "github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/webapp"
	"github.com/skygeario/skygear-server/pkg/core/async"
	pq2 "github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"net/http"
)

// Injectors from wire.go:

func newAuthorizeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	timeProvider := time.NewProvider()
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	urlProvider := &handler.URLProvider{
		Endpoints: endpointsProvider,
	}
	store := redis2.ProvideStore(context, tenantConfiguration, timeProvider)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	checker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, checker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	stateStoreImpl := &webapp.StateStoreImpl{
		Context: context,
	}
	webappURLProvider := &webapp.URLProvider{
		Endpoints: endpointsProvider,
		Anonymous: anonymousFlow,
		States:    stateStoreImpl,
	}
	scopesValidator := _wireScopesValidatorValue
	tokenGenerator := _wireTokenGeneratorValue
	authorizationHandler := handler.ProvideAuthorizationHandler(context, tenantConfiguration, factory, authorizationStore, grantStore, urlProvider, webappURLProvider, scopesValidator, tokenGenerator, timeProvider)
	httpHandler := provideAuthorizeHandler(factory, txContext, authorizationHandler)
	return httpHandler
}

var (
	_wireScopesValidatorValue = handler.ScopesValidator(oidc.ValidateScopes)
	_wireTokenGeneratorValue  = handler.TokenGenerator(oauth2.GenerateToken)
)

func newTokenHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	timeProvider := time.NewProvider()
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	store := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, store, authAccessEventProvider, tenantConfiguration)
	redisStore := redis2.ProvideStore(context, tenantConfiguration, timeProvider)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	checker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, checker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	urlprefixProvider := urlprefix.NewProvider(r)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration)
	interactionProvider := interaction.ProvideProvider(redisStore, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	httpHandler := provideTokenHandler(factory, txContext, tokenHandler)
	return httpHandler
}

func newRevokeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	revokeHandler := &handler.RevokeHandler{
		OfflineGrants: grantStore,
		AccessGrants:  grantStore,
	}
	httpHandler := provideRevokeHandler(factory, txContext, revokeHandler)
	return httpHandler
}

func newMetadataHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	metadataProvider := &oauth2.MetadataProvider{
		AuthorizeEndpoint: endpointsProvider,
		TokenEndpoint:     endpointsProvider,
		RevokeEndpoint:    endpointsProvider,
	}
	oidcMetadataProvider := &oidc.MetadataProvider{
		URLPrefix:          urlprefixProvider,
		JWKSEndpoint:       endpointsProvider,
		UserInfoEndpoint:   endpointsProvider,
		EndSessionEndpoint: endpointsProvider,
	}
	httpHandler := provideMetadataHandler(metadataProvider, oidcMetadataProvider)
	return httpHandler
}

func newJWKSHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	httpHandler := provideJWKSHandler(tenantConfiguration)
	return httpHandler
}

func newUserInfoHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	timeProvider := time.NewProvider()
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, store, userprofileStore, timeProvider)
	httpHandler := provideUserInfoHandler(factory, txContext, idTokenIssuer)
	return httpHandler
}

func newEndSessionHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	timeProvider := time.NewProvider()
	store := redis2.ProvideStore(context, tenantConfiguration, timeProvider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	checker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, checker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	stateStoreImpl := &webapp.StateStoreImpl{
		Context: context,
	}
	urlProvider := &webapp.URLProvider{
		Endpoints: endpointsProvider,
		Anonymous: anonymousFlow,
		States:    stateStoreImpl,
	}
	endSessionHandler := handler2.ProvideEndSessionHandler(tenantConfiguration, endpointsProvider, urlProvider, urlProvider)
	httpHandler := provideEndSessionHandler(factory, txContext, endSessionHandler)
	return httpHandler
}

func newChallengeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	validator := auth.ProvideValidator(m)
	context := auth.ProvideContext(r)
	timeProvider := time.NewProvider()
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	provider3 := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	challengeHandler := &ChallengeHandler{
		Validator:  validator,
		Challenges: provider3,
	}
	httpHandler := provideChallengeHandler(challengeHandler)
	return httpHandler
}

// wire.go:

func provideAuthorizeHandler(lf logging.Factory, tx db.TxContext, ah oauthAuthorizeHandler) http.Handler {
	h := &AuthorizeHandler{
		logger:       lf.NewLogger("oauth-authz-handler"),
		txContext:    tx,
		authzHandler: ah,
	}
	return h
}

func provideTokenHandler(lf logging.Factory, tx db.TxContext, th oauthTokenHandler) http.Handler {
	h := &TokenHandler{
		logger:       lf.NewLogger("oauth-token-handler"),
		txContext:    tx,
		tokenHandler: th,
	}
	return h
}

func provideRevokeHandler(lf logging.Factory, tx db.TxContext, rh oauthRevokeHandler) http.Handler {
	h := &RevokeHandler{
		logger:        lf.NewLogger("oauth-revoke-handler"),
		txContext:     tx,
		revokeHandler: rh,
	}
	return h
}

func provideMetadataHandler(oauth3 *oauth2.MetadataProvider, oidc2 *oidc.MetadataProvider) http.Handler {
	h := &MetadataHandler{
		metaProviders: []oauthMetadataProvider{oauth3, oidc2},
	}
	return h
}

func provideJWKSHandler(config2 *config.TenantConfiguration) http.Handler {
	h := &JWKSHandler{
		config: *config2.AppConfig.OIDC,
	}
	return h
}

func provideUserInfoHandler(lf logging.Factory, tx db.TxContext, uip oauthUserInfoProvider) http.Handler {
	h := &UserInfoHandler{
		logger:           lf.NewLogger("oauth-userinfo-handler"),
		txContext:        tx,
		userInfoProvider: uip,
	}
	return h
}

func provideEndSessionHandler(lf logging.Factory, tx db.TxContext, esh oidcEndSessionHandler) http.Handler {
	h := &EndSessionHandler{
		logger:            lf.NewLogger("oauth-end-session-handler"),
		txContext:         tx,
		endSessionHandler: esh,
	}
	return h
}

func provideChallengeHandler(h *ChallengeHandler) http.Handler {
	return h
}
