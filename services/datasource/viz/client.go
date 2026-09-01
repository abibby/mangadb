package viz

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/di"
	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/nulls"
	"github.com/ericchiang/css"
	"go.uber.org/ratelimit"
	"golang.org/x/net/html"
)

type Client struct {
	httpClient *http.Client

	limiter ratelimit.Limiter
}

type MangaIndex struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Chapters    []Chapter `json:"chapters"`
	Volumes     []Volume  `json:"volumes"`
}

type Chapter struct {
	ID       string     `json:"id"`
	Chapter  float32    `json:"chapter"`
	Volume   *nulls.Int `json:"volume"`
	VolumeID string     `json:"volume_id"`
	Link     string     `json:"link"`
}

type Volume struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Link   string `json:"link"`
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		limiter: ratelimit.New(20, ratelimit.Per(time.Minute)),
	}
}

func Register(ctx context.Context) {
	di.RegisterLazySingleton(ctx, func() (*Client, error) {
		return New(), nil
	})
}
func (c *Client) Series(ctx context.Context, tx database.DB, id string) error {
	u := "https://www.viz.com/shonenjump/chapters/" + id
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	document, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	index := MangaIndex{
		Chapters: []Chapter{},
		Volumes:  []Volume{},
	}

	title, err := css.Parse(`#series-intro h2`)
	if err != nil {
		return err
	}
	description, err := css.Parse(`#series-intro > div:first-child > div:last-child`)
	if err != nil {
		return err
	}
	chaptersNoVolume, err := css.Parse(`[data-sort-alpha-num] > [id^="ch-"]`)
	if err != nil {
		return err
	}
	chapters, err := css.Parse(`[id^="ch-"]`)
	if err != nil {
		return err
	}
	volumes, err := css.Parse(`.o_chapter-vol-container`)
	if err != nil {
		return err
	}
	volumeLink, err := css.Parse(`a[href^="/manga-books/manga"]`)
	if err != nil {
		return err
	}
	for _, ele := range title.Select(document) {
		index.Title = html.UnescapeString(ele.FirstChild.Data)
	}
	for _, ele := range description.Select(document) {
		index.Description = html.UnescapeString(ele.FirstChild.Data)
	}
	for _, ele := range chaptersNoVolume.Select(document) {
		c, err := chapter(ele, nil)
		if err != nil {
			return err
		}
		index.Chapters = append(index.Chapters, c)
	}

	for _, volume := range volumes.Select(document) {
		href := ""
		for _, link := range volumeLink.Select(volume) {
			for _, attr := range link.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
			}
		}

		parts := strings.SplitN(strings.TrimPrefix(href, "/manga-books/manga/"+id+"-volume-"), "-", 2)
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		num := nulls.NewInt(v)

		pathParts := strings.Split(href, "/")
		vol := Volume{
			ID:     pathParts[len(pathParts)-2],
			Number: num.Int(),
			Link:   "https://www.viz.com" + href,
		}
		index.Volumes = append(index.Volumes, vol)

		for _, ele := range chapters.Select(volume) {
			c, err := chapter(ele, &vol)
			if err != nil {
				return err
			}
			index.Chapters = append(index.Chapters, c)
		}
	}

	r, err := models.ApiResponseQuery(ctx).Where("url", "=", u).First(tx)
	if err != nil {
		return err
	}

	if r == nil {
		r = &models.APIResponse{
			Source:         "mangadex",
			SourceSeriesID: id,
			DataType:       "series",
			URL:            u,
			Page:           0,
		}
	}

	b, err := json.Marshal(index)
	if err != nil {
		return err
	}

	r.SyncJobID = ""
	r.RawPayload = b
	return model.SaveContext(ctx, tx, r)
}

func chapter(ele *html.Node, vol *Volume) (Chapter, error) {
	var err error
	var ch float64
	var href string

	for _, attr := range ele.Attr {
		if attr.Key == "name" {
			ch, err = strconv.ParseFloat(attr.Val, 32)
			if err != nil {
				return Chapter{}, err
			}
		}
		if attr.Key == "data-target-url" {
			href = attr.Val
		}
	}

	if strings.HasPrefix(href, "/") {
		href = "https://www.viz.com" + href
	} else if strings.HasPrefix(href, "javascript:tryReadChapter(") {
		href = "https://www.viz.com" + href[strings.Index(href, "'/")+1:len(href)-3]
	}

	u, err := url.Parse(href)
	if err != nil {
		return Chapter{}, err
	}

	parts := strings.Split(u.Path, "/")

	c := Chapter{
		ID:      parts[len(parts)-1],
		Chapter: float32(ch),
		Link:    href,
	}

	if vol != nil {
		c.Volume = nulls.NewInt(vol.Number)
		c.VolumeID = vol.ID
	}
	return c, nil
}
