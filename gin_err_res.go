package n_error

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrResponse struct {
	Code uint   `json:"code"`
	Msg  string `json:"msg"`
}

// ErrJSONResponse depends on:
//
//	"github.com/gin-contrib/requestid"
//	"github.com/gin-gonic/gin"
func ErrJSONResponse(ctx *gin.Context, err Tracer) {
	ctx.Error(err)

	//nowString := time.Now().Format(time.RFC3339)
	toDearEngineer := "請聯絡工程師查看 log"
	requestID := requestid.Get(ctx)

	if projectErr, ok := err.(*ProjectError); ok {
		// 錯誤有分客製訊息，跟真正的發生錯誤訊息
		// projectErr.Msg 是客製訊息
		actualErrMsg := projectErr.Msg

		if projectErr.Err != nil {
			actualErrMsg = projectErr.Err.Error()
		}

		fmt.Printf("%v：：錯誤代碼 %d：：路徑 %s：：錯誤內容 %s：：參數 %v：：EOF\n",
			requestID, projectErr.Code, projectErr.Tracer,
			actualErrMsg, projectErr.Params,
		)

		res := ErrResponse{
			Code: projectErr.Code,
			Msg:  projectErr.Msg,
		}

		// 針對特定錯誤，回覆 500 與固定錯誤訊息
		if SliceContain(projectErr.HttpStatus, []int{StatusDBError, http.StatusInternalServerError}) {
			res.Msg = toDearEngineer
			projectErr.HttpStatus = http.StatusInternalServerError
		}

		ctx.JSON(projectErr.HttpStatus, res)
	} else {
		fmt.Printf("發生非 project 錯誤: %v", err)
		ctx.String(http.StatusInternalServerError, toDearEngineer+" 2")
	}
}
