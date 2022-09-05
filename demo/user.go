package demo

import (
	"fmt"
	"webv5"
)

func Main(c *webv5.Context) {
	fmt.Fprintf(c.W, "this is main")
}

func Mainn(c *webv5.Context) {
	fmt.Fprintf(c.W, "this is mainn")
}

type signUpReq struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ConfirmrPassword string `json:"confirm_password"`
}

type commonResponse struct {
	BizCode int
	Msg     string
	Data    int
}

// 不再由用户手动创建context
func SignUp(c *webv5.Context) {
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		_ = c.BadRequestJson(&commonResponse{
			BizCode: 4,
			Msg:     fmt.Sprintf("invalid request: %v", err),
		})
		return
	}
	_ = c.BadRequestJson(&commonResponse{
		Data: 123,
	})
	fmt.Fprintf(c.W, "\ndata is: %v", *req)
}
