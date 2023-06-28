package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Env = "release"

type ErrResponseHandler func(ctx *gin.Context, env string, err error, statusCode int)
type SuccessResponseHandler func(ctx *gin.Context, data any, statusCode int)

var errResponseHandler ErrResponseHandler = func(ctx *gin.Context, env string, err error, statusCode int) {
	ctx.JSON(statusCode, map[string]interface{}{
		"error": err.Error(),
	})
}

var successResponseHandler SuccessResponseHandler = func(ctx *gin.Context, data any, statusCode int) {
	ctx.JSON(statusCode, map[string]interface{}{
		"data": data,
		"code": statusCode,
	})
}

func SetErrResponseHandler(handler ErrResponseHandler) {
	errResponseHandler = handler
}

func SetSuccessResponseHandler(handler SuccessResponseHandler) {
	successResponseHandler = handler
}

func GetErrResponseHandler() ErrResponseHandler {
	return errResponseHandler
}

func GetSuccessResponseHandler() SuccessResponseHandler {
	return successResponseHandler
}

type Response struct {
}

func (response Response) JsonSuccessResponse(ctx *gin.Context) {
	response.JsonResponseWithoutError(ctx, "success")
}

func (response Response) JsonResponseWithoutError(ctx *gin.Context, data any) {
	response.JsonResponse(ctx, data, nil, http.StatusOK)
}

func (response Response) JsonResponseWithServerError(ctx *gin.Context, err error) {
	response.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
}

func (response Response) JsonResponseWithError(ctx *gin.Context, err error, statusCode int) {
	response.JsonResponse(ctx, "", err, statusCode)
}

func (response Response) JsonResponse(ctx *gin.Context, data any, err error, statusCode int) {
	ctx.Set("response_data", data)
	ctx.Set("response_err", err)
	ctx.Set("response_code", statusCode)

	if err != nil {
		GetErrResponseHandler()(ctx, Env, err, statusCode)
		return
	}

	GetSuccessResponseHandler()(ctx, data, statusCode)
}
