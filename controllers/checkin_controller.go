package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"encoding/json"
	"net/http"
)

// Checkin 创建打卡任务
func Checkin(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	var checkin models.Checkin

	// 解析请求体中的数据
	if err := json.NewDecoder(r.Body).Decode(&checkin); err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 调用服务层来创建打卡记录
	ctx := context.Background()
	createdCheckin, err := service.CreateCheckinService(ctx, checkin)
	if err != nil {
		// 返回服务层错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功创建打卡任务，返回 201 状态码和响应信息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkin created successfully", Data: createdCheckin})
}

// GetCheckinRecords 获取用户的打卡记录
func GetCheckinRecords(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "User ID is required"})
		return
	}

	// 调用服务层来获取用户的打卡记录
	ctx := context.Background()
	checkins, err := service.GetCheckinRecords(ctx, userID)
	if err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回用户的打卡记录
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkin records retrieved successfully", Data: checkins})
}
