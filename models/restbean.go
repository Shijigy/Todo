package models

type RestBean struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

func NewRestBean(status int, success bool, message interface{}) *RestBean {
	return &RestBean{
		Status:  status,
		Success: success,
		Message: message,
	}
}

func SuccessRestBean() *RestBean {
	return NewRestBean(200, true, nil)
}

func SuccessRestBeanWithData(data interface{}) *RestBean {
	return NewRestBean(200, true, data)
}

func FailureRestBean(status int) *RestBean {
	return NewRestBean(status, false, nil)
}

func FailureRestBeanWithData(status int, data interface{}) *RestBean {
	return NewRestBean(status, false, data)
}
