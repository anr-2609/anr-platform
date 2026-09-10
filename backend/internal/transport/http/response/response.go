package response

import (
	"encoding/json"
	"net/http"
)

// Response là format JSON response chuẩn cho toàn bộ API của ANR Platform.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError định nghĩa thông tin lỗi có cấu trúc trả về client.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON ghi response thành công ra ResponseWriter dưới dạng JSON.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// Error ghi response lỗi ra ResponseWriter dưới dạng JSON chuẩn.
func Error(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
