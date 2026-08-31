package mangaplus

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"

	"abibby.com/salusa/di"
	"github.com/abibby/icbmdb/services/datasource/mangaplus/mpproto"
)

type Client struct {
	proto mpproto.Client
}

func NewClient(c *http.Client) *Client {
	return &Client{
		proto: *mpproto.NewClient(c),
	}
}

func Register(ctx context.Context) {
	di.RegisterLazySingleton(ctx, func() (*Client, error) {
		return NewClient(&http.Client{
			Timeout: time.Second * 10,
		}), nil
	})
}
func (c *Client) TitleDetailsV3(id string) (*mpproto.TitleDetailView, error) {

	result, err := c.proto.Get("https://jumpg-webapi.tokyo-cdn.com/api/title_detailV3?title_id=%s", id)
	if err != nil {
		return nil, err
	}
	return result.GetTitleDetailView(), nil
}

func (c *Client) Register(deviceID string) (*mpproto.RegistrationData, error) {

	deviceToken := md5.New()
	deviceToken.Write([]byte(deviceID))
	deviceTokenStr := hex.EncodeToString(deviceToken.Sum([]byte{}))

	securityKey := md5.New()
	securityKey.Write([]byte(deviceTokenStr + "4Kin9vGg"))
	securityKeyStr := hex.EncodeToString(securityKey.Sum([]byte{}))

	body := url.Values{}
	body.Add("device_token", deviceTokenStr)
	body.Add("security_key", securityKeyStr)
	result, err := c.proto.Put(body, "https://jumpg-webapi.tokyo-cdn.com/api/register")
	if err != nil {
		return nil, err
	}
	return result.GetRegisterationData(), nil
}
