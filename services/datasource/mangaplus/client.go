package mangaplus

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"abibby.com/mangadb/services/datasource/mangaplus/mpproto"
	"gosalusa.com/database"
	"gosalusa.com/di"
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

// Series implements [datasource.Datasource].
func (c *Client) Series(ctx context.Context, tx database.DB, id string) error {

	u := fmt.Sprintf("https://jumpg-webapi.tokyo-cdn.com/api/title_detailV3?title_id=%s", id)
	result, err := c.proto.Get(u)
	if err != nil {
		return err
	}

	b, err := json.Marshal(result.GetTitleDetailView())
	if err != nil {
		return err
	}
	return models.ApiResponseCreateOrUpdate(ctx, tx, &models.APIResponse{
		Source:         datasource.MangaplusSource,
		SourceSeriesID: id,
		DataType:       "series",
		URL:            u,
		Page:           0,
		SyncJobID:      "",
		RawPayload:     b,
	})
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
