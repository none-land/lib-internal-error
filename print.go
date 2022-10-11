package n_error

import "fmt"

func Print(requestID string, nErr Tracer) {
	if projectErr, ok := nErr.(*ProjectError); ok {
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
	}
}
