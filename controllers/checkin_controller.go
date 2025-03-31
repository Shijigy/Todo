package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Checkin 处理打卡请求
func Checkin(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	var checkin models.Checkin

	// 解析前端传来的数据
	if err := json.NewDecoder(r.Body).Decode(&checkin); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 解析 start_date 和 end_date
	startDate, err := time.Parse("2006-01-02", checkin.StartDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", checkin.EndDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid end_date format"})
		return
	}

	// 初始化打卡次数和状态
	checkinCount := make(map[string]int)

	// 遍历日期范围，初始化每天的打卡次数为 0
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		checkinCount[dateStr] = 0 // 初始化每个日期的打卡次数为 0
	}

	// 将 checkinCount 序列化为 JSON 字符串
	checkinCountJSON, err := json.Marshal(checkinCount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to marshal checkin count"})
		return
	}

	// 设置 checkin 的 CheckinCount 为 JSON 字符串
	checkin.CheckinCount = checkinCountJSON

	// 调用服务层来创建打卡记录
	ctx := context.Background()
	createdCheckin, err := service.CreateCheckinService(ctx, checkin)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功创建的打卡记录
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// 使用 decodedCount 代替 CheckinCount
	createdCheckin.Checkin.CheckinCount = []byte(fmt.Sprintf("%s", createdCheckin.DecodedCheckin))

	// 返回响应
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkin created successfully",
		Data:    createdCheckin,
	})
}

// GetCheckinsByUserAndDate 获取指定用户指定日期的打卡记录
/*func GetCheckinsByUserAndDate(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 从请求参数中获取 user_id 和 date
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	// 校验参数
	if userID == "" || date == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Missing user_id or date"})
		return
	}

	// 调用服务层获取指定用户指定日期的打卡记录
	checkins, err := service.GetCheckinsByUserIDAndDateService(r.Context(), userID, date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回打卡记录
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Checkins retrieved successfully", Data: checkins})
}
*/
// GetCheckinsByUserAndDate 获取指定用户指定日期的打卡记录
func GetCheckinsByUserAndDate(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 从请求参数中获取 user_id 和 date
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	// 校验参数
	if userID == "" || date == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Missing user_id or date"})
		return
	}

	// 调用服务层获取指定用户指定日期的打卡记录
	checkins, seenIDs, decodedCounts, err := service.GetCheckinsByUserIDAndDateService(r.Context(), userID, date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 打印日志
	fmt.Printf("Returned checkins: %+v\n", checkins)
	fmt.Printf("Returned seenIDs: %+v\n", seenIDs)
	fmt.Printf("Returned decodedCounts: %+v\n", decodedCounts)

	// 返回打卡记录，包含 seenIDs 和 decodedCounts
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkins retrieved successfully",
		Data: map[string]interface{}{
			"checkins":      checkins,
			"seenIDs":       seenIDs,
			"decodedCounts": decodedCounts,
		},
	})
}

func GetCheckinRecordByID(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	checkinIDStr := r.URL.Query().Get("checkin_id")
	fmt.Println("接收到的 checkin_id:", checkinIDStr)

	if checkinIDStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "checkin_id 是必填项"})
		return
	}

	checkinID, err := strconv.Atoi(checkinIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "checkin_id 格式无效"})
		return
	}

	// 调试日志，检查流程
	fmt.Println("调用服务获取打卡记录，ID:", checkinID)

	checkinWithDecoded, err := service.GetCheckinByIDService(r.Context(), checkinID)
	if err != nil {
		fmt.Println("获取打卡记录时出错:", err) // 输出具体错误
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "获取打卡记录失败"})
		return
	}

	// 检查是否有打卡记录
	if checkinWithDecoded == nil || checkinWithDecoded.Checkin == nil {
		fmt.Println("没有找到 ID 对应的打卡记录:", checkinID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Error: "没有找到打卡记录"})
		return
	}

	// 返回解码后的 checkin_count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: "打卡记录成功获取",
		Data: map[string]interface{}{
			"checkin":               checkinWithDecoded.Checkin,
			"decoded_checkin_count": checkinWithDecoded.DecodedCheckin,
		},
	})
}

func IncrementCheckinCount(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 定义请求体结构
	var requestData struct {
		CheckinID int    `json:"checkin_id"`
		Date      string `json:"date"`
	}

	// 从请求体中读取 JSON 数据
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		// 如果解析请求体失败，返回错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid request format"})
		return
	}

	// 校验参数
	if requestData.CheckinID == 0 || requestData.Date == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Missing checkin_id or date"})
		return
	}

	// 调用服务层更新打卡次数并获取更新后的打卡信息
	checkin, err := service.IncrementCheckinCountService(r.Context(), requestData.CheckinID, requestData.Date)
	if err != nil {
		// 如果服务层返回错误，判断错误类型
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == fmt.Sprintf("checkin already completed for the date %s", requestData.Date) {
			// 如果打卡已完成，返回特定错误
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			json.NewEncoder(w).Encode(models.Response{Error: "Check-in already completed for this date"})
		} else {
			// 如果其他错误，返回内部服务器错误
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Response{Error: "Failed to increment checkin count: " + err.Error()})
		}
		return
	}

	// 返回成功响应，包含更新后的打卡信息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkin count incremented successfully",
		Data:    checkin, // 返回更新后的打卡信息
	})
}

// CheckinCompleted 判断打卡是否完成
func CheckinCompleted(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 解析请求体
	var req models.CheckinCompletionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "无效的请求参数"})
		return
	}

	// 调用服务层判断是否完成打卡
	completed, err := service.CheckIfCheckinCompleted(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: fmt.Sprintf("服务层出错：%v", err)})
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	if completed {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Response{Message: "打卡已完成"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Response{Message: "打卡未完成"})
	}
}

// UpdateCheckin 修改打卡信息
func UpdateCheckin(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	var checkin models.Checkin

	// 解析前端传来的 JSON 数据
	if err := json.NewDecoder(r.Body).Decode(&checkin); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 获取 ID（直接从请求体中获取）
	if checkin.ID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Checkin ID is required"})
		return
	}

	// 解析 start_date 和 end_date
	startDate, err := time.Parse("2006-01-02", checkin.StartDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", checkin.EndDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid end_date format"})
		return
	}

	// 调用服务层来更新打卡记录
	ctx := context.Background()
	updatedCheckin, err := service.UpdateCheckinService(ctx, checkin.ID, checkin, startDate, endDate) // 传递 id（int）
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功更新的打卡记录
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 返回响应
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkin updated successfully",
		Data:    updatedCheckin,
	})
}

// MarkCheckinCompleted 处理标记某天某打卡完成的请求
func MarkCheckinCompleted(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	// 解析前端传来的数据
	var request struct {
		CheckinID int    `json:"checkin_id"`
		Date      string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 解析 date 为 time 类型
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid date format"})
		return
	}

	// 调用服务层获取打卡记录
	ctx := r.Context()
	checkin, err := service.GetCheckinByIDService(ctx, request.CheckinID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 检查日期是否存在于打卡记录中
	if checkin.DecodedCheckin == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "No checkin data found for the specified checkin_id"})
		return
	}

	// 将指定日期的打卡次数设置为目标打卡次数
	dateStr := date.Format("2006-01-02")
	if _, exists := checkin.DecodedCheckin[dateStr]; exists {
		checkin.DecodedCheckin[dateStr] = checkin.Checkin.TargetCheckinCount
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Date not found in checkin record"})
		return
	}

	// 将更新后的数据转换为 JSON
	updatedCheckinCount, err := json.Marshal(checkin.DecodedCheckin)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to marshal updated checkin count"})
		return
	}

	// 更新数据库中的打卡记录
	err = service.UpdateCountService(ctx, request.CheckinID, *checkin.Checkin, checkin.Checkin.StartDate, checkin.Checkin.EndDate, updatedCheckinCount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkin marked as completed for the specified date",
		Data:    checkin,
	})
}

// DeleteCheckin 删除打卡
func DeleteCheckin(c *gin.Context, service services.CheckinService) {
	var body struct {
		CheckinID string `json:"checkin_id"` // 接受字符串类型
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 转换 checkin_id 字符串为整数
	checkinID, err := strconv.Atoi(body.CheckinID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkin_id"})
		return
	}

	err = service.DeleteCheckinService(c, checkinID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkin deleted successfully"})
}

// ResetCheckinCount 处理更新打卡次数的请求
func ResetCheckinCount(w http.ResponseWriter, r *http.Request, service services.CheckinService) {
	var req struct {
		CheckinID  int    `json:"checkin_id"`
		Date       string `json:"date"`
		ResetCount int    `json:"reset_count"`
	}

	// 解析前端传来的数据
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 获取当前日期并验证格式
	_, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid date format"})
		return
	}

	// 获取打卡记录
	ctx := context.Background()
	checkin, err := service.GetCheckinByIDService(ctx, req.CheckinID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error fetching checkin data"})
		return
	}

	// 检查是否超过目标打卡次数
	if req.ResetCount > checkin.Checkin.TargetCheckinCount {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Reset count exceeds target checkin count"})
		return
	}

	// 解码 CheckinCount 为 map
	var checkinCount map[string]int
	err = json.Unmarshal(checkin.Checkin.CheckinCount, &checkinCount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to decode checkin count"})
		return
	}

	// 更新指定日期的打卡次数
	checkinCount[req.Date] = req.ResetCount

	// 将更新后的打卡次数重新编码为 JSON
	updatedCheckinCountJSON, err := json.Marshal(checkinCount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to marshal updated checkin count"})
		return
	}

	// 更新数据库中的打卡次数
	checkin.Checkin.CheckinCount = updatedCheckinCountJSON
	err = service.ResetCheckinCountService(ctx, checkin.Checkin.ID, checkin.Checkin.CheckinCount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to update checkin count"})
		return
	}

	// 返回更新后的打卡记录
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: "Checkin count updated successfully",
		Data:    checkin,
	})
}
