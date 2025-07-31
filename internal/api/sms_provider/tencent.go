package sms_provider

import (
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type TencentSMSProvider struct {
	client *sms.Client
}

func NewTencentSMSProvider() *TencentSMSProvider {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 10 // 请求超时时间，单位为秒(默认60秒)
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	return &TencentSMSProvider{client: client}
}
func (p *TencentSMSProvider) SendMessage(phone, message, channel, otp string) (string, error) {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr("1400787878")
	request.SignName = common.StringPtr("腾讯云")
	request.TemplateId = common.StringPtr("449739")
	request.TemplateParamSet = common.StringPtrs([]string{"1234"})
	request.PhoneNumberSet = common.StringPtrs([]string{"+8613711112222"})
	request.SessionContext = common.StringPtr("")
	request.ExtendCode = common.StringPtr("")
	request.SenderId = common.StringPtr("")
	response, err := p.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	return *response.Response.SendStatusSet[0].SerialNo, nil
}

func (p *TencentSMSProvider) VerifyOTP(phone, token string) error {
	panic("")
}
