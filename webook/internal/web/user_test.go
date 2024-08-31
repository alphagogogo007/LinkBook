package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T){
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode int
		wantBody string
	}{
		{
			name: "SuccessSignup",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:   "123@qq.com",
					Password: "12345678",
				}).Return(nil)
				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"12345678",
					 "confirmPassword":"12345678" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "hello, successfully signing up",

		},
		{
			name: "BindError",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"12345678
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusBadRequest,
			wantBody: "",

		},
		{
			name: "WrongEmail",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq",
					 "password":"12345678",
					 "confirmPassword":"12345678" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "illegal email name",

		},
		{
			name: "wrong confirm password",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"12345678",
					 "confirmPassword":"87654321" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "The passwords entered do not match",

		},
		{
			name: "wrong password format",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"123456",
					 "confirmPassword":"123456" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "The password format is incorrect; it must be at least eight characters long.",

		},
		{
			name: "system error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:   "123@qq.com",
					Password: "12345678",
				}).Return(errors.New("db error"))
				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"12345678",
					 "confirmPassword":"12345678" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "system error",

		},
		{
			name: "email conflict",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService){
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:   "123@qq.com",
					Password: "12345678",
				}).Return(service.ErrDuplicateEmail)
				return userSvc,nil
			},
			reqBuilder: func(t *testing.T) *http.Request{
				req,err := http.NewRequest(http.MethodPost,
					 "/users/signup",bytes.NewReader([]byte(`{
					 "email":"123@qq.com",
					 "password":"12345678",
					 "confirmPassword":"12345678" 
					 }`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			}, 
			wantCode: http.StatusOK,
			wantBody: "Email conflict, please use a different one.",

		},
		


	}

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			// before t.Run finish, it will execute finish
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterRoutes(server)
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			//assert
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())

		})
	}
}


func TestMOck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := svcmocks.NewMockUserService(ctrl)
	userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(nil)

	err:=userSvc.SignUp(context.Background(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	})
	t.Log(err)
}
