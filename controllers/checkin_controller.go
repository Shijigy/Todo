package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

// GetCheckinRecordByUserID 获取用户的所有打卡记录
func GetCheckinRecordByUserID(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	userID := r.URL.Query().Get("user_id")

	// 如果没有提供 user_id，返回错误
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "User ID is required"})
		return
	}

	// 调用服务层获取用户的所有打卡记录
	ctx := context.Background()
	checkins, err := service.GetCheckinRecordsByUserID(ctx, userID)
	if err != nil {
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

// MarkCheckinComplete 标记打卡任务完成
func MarkCheckinComplete(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	checkinID := r.URL.Query().Get("checkin_id")
	if checkinID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Checkin ID is required"})
		return
	}

	ctx := context.Background()
	updatedCheckin, err := service.MarkCheckinCompleteService(ctx, checkinID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkin marked as complete", Data: updatedCheckin})
}

// UpdateCheckinCount 更新打卡次数
func UpdateCheckinCount(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 获取 request 中的参数
	checkinID := r.URL.Query().Get("checkin_id")
	increment := r.URL.Query().Get("increment")

	// 校验参数
	if checkinID == "" || increment == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Checkin ID and increment are required"})
		return
	}

	// 转换增量为整数
	incrementInt, err := strconv.Atoi(increment)
	if err != nil || incrementInt <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid increment value"})
		return
	}

	// 调用服务层进行打卡次数更新
	ctx := context.Background()
	updatedCheckin, err := service.UpdateCheckinCountService(ctx, checkinID, incrementInt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkin count updated successfully", Data: updatedCheckin})
}

// DeleteCheckin 删除打卡任务
func DeleteCheckin(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	checkinID := r.URL.Query().Get("checkin_id")
	if checkinID == "" {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Checkin ID is required"})
		return
	}

	// 调用服务层删除打卡任务
	ctx := context.Background()
	err := service.DeleteCheckinService(ctx, checkinID)
	if err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功删除的响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkin deleted successfully"})
}
