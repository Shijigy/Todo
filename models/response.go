package models

// Response 用于统一的 JSON 响应结构
type Response struct {
	Message string      `json:"message"`         // 响应信息
	Data    interface{} `json:"data,omitempty"`  // 可选的数据
	Error   string      `json:"error,omitempty"` // 可选的错误信息
	Token   string      `json:"token,omitempty"` // 可选的 Token，注册成功时返回
	Status  string      `json:"status"`
}
