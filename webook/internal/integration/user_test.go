package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init(){
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandler_SendSMSCode(t *testing.T) {

	rdb := startup.InitRedis()
	server := startup.InitWebServer()

	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		phone    string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "send success",
			before: func(t *testing.T){},
			after: func(t *testing.T){
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.Get(ctx,  key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code)>0)
				duration, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, duration>time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t,  err)
			},
			phone: "15212341234",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "Successfully send the code",
			},
		},
		{
			name: "empty phone",
			before: func(t *testing.T){},
			after: func(t *testing.T){

			},
			phone: "",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "Please input phone number",
			},
		},
		{
			name: "send too many",
			before: func(t *testing.T){
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				err := rdb.Set(ctx, key, "123456", time.Minute*10).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T){
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.GetDel(ctx,  key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)

			},
			phone: "15212341234",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "Send too many",
			},
		},
		{
			name: "system error",
			before: func(t *testing.T){
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T){
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.GetDel(ctx,  key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)

			},
			phone: "15212341234",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "System error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

		
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			//assert
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)

		})
	}
}
