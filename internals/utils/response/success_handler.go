package response

import (
	"net/http"

	"wakuwaku_nihongo/internals/abstraction"

	"github.com/labstack/echo/v4"
)

type successConstant struct {
	OK                    Success
	InternalVALIDConflict Success // use in internal transfer va if callback id conflict
}

var SuccessConstant successConstant = successConstant{
	OK: Success{
		Response: successResponse{
			Meta: Meta{
				Success:    true,
				Message:    "Request successfully proceed",
				StatusCode: 20001,
			},
			Data: nil,
		},
		Code: http.StatusOK,
	},
	InternalVALIDConflict: Success{
		Response: successResponse{
			Meta: Meta{
				Success:    true,
				Message:    "id and va already create",
				StatusCode: 20002,
			},
			Data: nil,
		},
		Code: http.StatusOK,
	},
}

type successResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Success struct {
	Response successResponse `json:"response"`
	Code     int             `json:"code"`
}

func SuccessBuilder(res *Success, data interface{}) *Success {
	res.Response.Data = data
	return res
}

func SuccessResponse(data interface{}) *Success {
	return SuccessBuilder(&SuccessConstant.OK, data)
}

func (s *Success) Send(c echo.Context) error {
	return c.JSON(s.Code, s.Response)
}

type SuccessResponseWithInfo struct {
	Meta Meta                        `json:"meta"`
	Data interface{}                 `json:"data"`
	Info *abstraction.PaginationInfo `json:"info"`
}

func SuccessResponseInfo(data interface{}, info *abstraction.PaginationInfo) (res *SuccessResponseWithInfo) {
	res = &SuccessResponseWithInfo{
		Meta: SuccessConstant.OK.Response.Meta,
		Data: data,
		Info: info,
	}

	return
}

func (s *SuccessResponseWithInfo) Send(c echo.Context) error {
	return c.JSON(200, s)
}
