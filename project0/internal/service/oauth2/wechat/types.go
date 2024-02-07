package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"project0/internal/domain"
	"project0/pkg/loggerDefine"
)

type Service interface {
	// 构造url
	AuthURL(ctx context.Context, state string) (string, error)
	// verify code
	Verify(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	// 基本不变的字段也就是可以当常量了放这了,这种垃圾，低阶解释就写一遍，后续提交的时候抹掉他
	appID string
	//
	appSecret string

	client *http.Client

	logger loggerDefine.LoggerV1
}

// 这些要一个个深刻联系在一起。  怎么使用。 增强联系。
type Result struct {
	//接口调用凭证
	AccessToken string `json:"access_token"`
	//access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	//用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	//当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionID string `json:"unionid"`
	//错误返回

	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func NewService(appID, appSecret string,l loggerDefine.LoggerV1) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	    logger: l,}
}

func (s *service) Verify(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokeyUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	// 编码构造http请求,用的是getmethod，body 传nil
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokeyUrl, nil)
	if err != nil {
		log.Println("verify wechat access err : ", err)
		//goland的强大   直接ret就补全可选因为err返回的非err参数的空结构体及err
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		log.Println("verify wechat code failed ", err)
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Println("decode response body failed ", err)
		return domain.WechatInfo{}, err
	}

	if res.Errcode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败, %s,%s", res.Errcode, res.Errmsg)
	}

	return domain.WechatInfo{
		res.OpenID,
		res.UnionID,
	}, nil

}
func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`


	//state := uuid.New()

	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil

}
