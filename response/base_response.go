package response

type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Result  any    `json:"result,omitempty"`
}

type BasePageResult struct {
	Rows  any   `json:"rows"`
	Total int64 `json:"total"`
}

func ReturnErrorResponse(code int32, message string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Success: false,
	}
}

func ReturnSuccessPageResultResponse(code int32, message string, rows any, total int64) *BaseResponse {
	basePageResult := &BasePageResult{
		Total: total,
		Rows:  rows,
	}
	return &BaseResponse{
		Code:    code,
		Message: message,
		Success: true,
		Result:  basePageResult,
	}
}

func ReturnSuccessResultResponse(code int32, message string, result any) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Success: true,
		Result:  result,
	}
}

func ReturnSuccessNoResultResponse() *BaseResponse {
	return &BaseResponse{
		Code:    0,
		Message: "success",
		Success: true,
	}
}
