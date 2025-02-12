package config

// 定义应用的常量

// HTTP 状态码常量
const (
	OK                  = 200
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	InternalServerError = 500
)

// 环境常量
const (
	Development = "development"
	Production  = "production"
)

// 日志常量
const (
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarning = "WARNING"
	LogLevelError   = "ERROR"
)

// 路由常量
const (
	RouteLogin  = "/login"
	RouteLogout = "/logout"
	RouteHome   = "/home"
	RouteUser   = "/user"
)

// 默认值常量
const (
	DefaultPageSize = 10
	DefaultLanguage = "en"
)
