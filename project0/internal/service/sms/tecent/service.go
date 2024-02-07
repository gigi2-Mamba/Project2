package tecent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	tecnentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

/**
保持依赖注入，不关心短信client的生成
*/

type Service struct {
	client   *tecnentsms.Client
	appId    *string
	signName *string
}

// 腾讯云设计并不好
func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me
	request := tecnentsms.NewSendSmsRequest()
	//request.service.UserServiceontext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	for _, statuspPtr := range response.Response.SendStatusSet {
		if statuspPtr != nil {
			// 基本不可能进来这里
			continue
		}
		status := *statuspPtr
		if *status.Code != "Ok" || status.Code == nil {
			// 发送失败
			return fmt.Errorf("发送短信失败 code: %s, msg： %s", *status.Code, *status.Message)
		}
	}

	return nil
}

// ekit里面的方法值得研究啊
func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data, func(idx int, src string) *string {
		return &src
	})

}

func NewService(c *tecnentsms.Client, appId string, signName string) *Service {

	return &Service{
		client:   c,
		appId:    &appId,
		signName: &signName,
	}
}
