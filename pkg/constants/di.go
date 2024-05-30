package constants

const (
	ConfigName             = "Config"
	LoggerName             = "Logger"
	ValidatorName          = "Validator"
	ExecutorFactoryName    = "ExecutorFactory"
	ExecutorCtxFactoryName = "ExecutorCtxFactory"
	PgSQLName              = "PgSQL"
	MongoClientName        = "MongoClient"
	MongoDatabaseName      = "MongoDatabase"
	RedisName              = "Redis"
	HTTPServerName         = "HTTP"

	JWTAuthorizerName = "JWTAuthorizer"

	LogMiddlewareName = "LogMiddleware"

	NonZeroCustomRuleName   = "NonZeroCustomRule"
	BeforeNowCustomRuleName = "BeforeNowCustomRule"
	AfterNowCustomRuleName  = "AfterNowCustomRule"

	JWTMiddlewareFactoryName = "JWTMiddlewareFactory"
	HTTPErrorHandlerName     = "ErrorHTTPHandler"
)
