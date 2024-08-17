package httputils

import (
	"encoding/json"
	"fmt"
	"go-chi-boilerplate/utils"
	"go-chi-boilerplate/utils/constants"
	"math"
	"net/http"
)

type (
	Errors   map[string]string
	BaseMeta struct {
		Page      int `json:"page"`
		Limit     int `json:"limit"`
		TotalData int `json:"total_data"`
		TotalPage int `json:"total_page"`
	}

	// BaseResponse is the base response
	BaseResponse struct {
		Code       string      `json:"code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		Meta       *BaseMeta   `json:"meta"`
		Errors     []Errors    `json:"errors"`
		ServerTime int64       `json:"server_time"`
	}
)

// MapBaseResponse map response
func MapBaseResponse(w http.ResponseWriter, r *http.Request, message string, data interface{}, meta *BaseMeta, err error, errors []Errors) {
	// Check Request ID
	requestID := r.Header.Get("request-id")
	if requestID != "" {
		bodyByte, _ := json.Marshal(data)
		fmt.Println("[RESPONSE: ", r.URL.String(), "] REQUEST_ID: ", requestID, " BODY:", string(bodyByte))
	}

	statusCode, code := utils.GetStatusCode(err)

	// Payload Response
	payload := BaseResponse{
		Code:       code,
		Message:    message,
		Data:       data,
		Errors:     errors,
		ServerTime: utils.TimeNow().Unix(),
		Meta:       meta,
	}

	// Marshal json response
	jsonResponse, _ := json.MarshalIndent(payload, "", "	")

	// Write Response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Date", utils.TimeNow().Format(constants.FORMAT_DATETIME_TEXT))
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func SetBaseMeta(page int, limit int, totalData int) BaseMeta {
	totalPage := float64(totalData) / float64(limit)
	return BaseMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: int(math.Ceil(totalPage)),
	}
}
