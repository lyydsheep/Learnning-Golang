package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service"
	svcmocks "github.com/lyydsheep/Learnning-Golang/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestEncrypt(t *testing.T) {
//	password := "123456"
//	//对密码进行加密，并获取加密后的结果
//	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		t.Fatal(err)
//	}
//	//比较密码是否一致
//	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
//	assert.NoError(t, err)
//}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		input        string
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "场景——注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				usvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return usvc, nil
			},
			input:        `{"email":"123@qq.com","password":"123o456","confirmPassword":"123o456"}`,
			expectedCode: 200,
			expectedMsg:  "注册成功",
		},
		{
			name: "场景——Bind失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				return usvc, nil
			},
			input:        `"email":"123@qq.com","password":"123o456","confirmPassword":"123o456"`,
			expectedCode: 400,
			expectedMsg:  "",
		},
		{
			name: "场景——两次输入密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				return usvc, nil
			},
			input:        `{"email":"123@qq.com","password":"123o4156","confirmPassword":"123o456"}`,
			expectedCode: 200,
			expectedMsg:  "两次输入密码不一致",
		},
		{
			name: "场景——邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				return usvc, nil
			},
			input:        `{"email":"123@qq","password":"123o456","confirmPassword":"123o456"}`,
			expectedCode: 200,
			expectedMsg:  "邮箱格式不对",
		},
		{
			name: "场景——密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				return usvc, nil
			},
			input:        `{"email":"123@qq.com","password":"123@o456","confirmPassword":"123@o456"}`,
			expectedCode: 200,
			expectedMsg:  "只能由字母、数字组成，1-9位",
		},
		{
			name: "场景——邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				usvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return usvc, nil
			},
			input:        `{"email":"123@qq.com","password":"123o456","confirmPassword":"123o456"}`,
			expectedCode: 200,
			expectedMsg:  "邮箱冲突",
		},
		{
			name: "场景——系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usvc := svcmocks.NewMockUserService(ctrl)
				usvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("随意一个错误"))
				return usvc, nil
			},
			input:        `{"email":"123@qq.com","password":"123o456","confirmPassword":"123o456"}`,
			expectedCode: 200,
			expectedMsg:  "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usvc, csvc := tc.mock(ctrl)

			u := NewUserHandler(usvc, csvc)
			gin.SetMode(gin.ReleaseMode)
			server := gin.Default()
			u.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)
			assert.Equal(t, tc.expectedMsg, resp.Body.String())
		})
	}
}

func TestUserHandler_Profile(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func(ctrl gomock.Controller) (service.UserService, service.CodeService)
		input        string
		expectCode   int
		expectResult Result
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
