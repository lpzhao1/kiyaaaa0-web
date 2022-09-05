package demo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webv0"
)

func Main(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is main")
}

type signUpReq struct {
	Email            string `json:"email`
	Password         string `json:"password`
	ConfirmrPassword string `json:"confirm_password"`
}

type commonResponse struct {
	BizCode int
	Msg     string
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	c := webv0.NewContext(w, r)
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		resp := &commonResponse{
			BizCode: 4, //假设这个代码代表输入参数错误
			Msg:     fmt.Sprintf("invalid request: %v", err),
		}
		respBytes, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(respBytes))
		return
	}
	fmt.Fprintf(w, "invalid request: %v", err)

}
