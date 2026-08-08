package mangaplus

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/icbmdb/services/mangaplus"
)

type MangaPlus struct {
	client *mangaplus.Client
}

var _ datasource.Datasource = (*MangaPlus)(nil)

func New() *MangaPlus {
	return &MangaPlus{
		client: mangaplus.NewClient(http.DefaultClient),
	}
}

// Quality implements [datasource.Datasource].
func (m *MangaPlus) Quality() int {
	return 200
}

// ID implements [datasource.Datasource].
func (m *MangaPlus) ID(s *models.DatasourceIDs) string {
	return s.MangaPlusID
}

// UpdateTimestamp implements [datasource.Datasource].
func (m *MangaPlus) UpdateTimestamp(s *models.DatasourceIDs) {
	s.MangaPlusUpdated = time.Now()
}

// Series implements [datasource.Datasource].
func (m *MangaPlus) Series(id string) (datasource.SeriesApplier, error) {
	titleDetail, err := m.client.TitleDetailsV3(id)
	if err != nil {
		return nil, err
	}

	return &Series{
		client:      m.client,
		titleDetail: titleDetail,
	}, nil
}

// ExtractSeriesID implements [datasource.Datasource].
func (m *MangaPlus) ApplySeriesID(s *models.Series, uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}
	if u.Hostname() != "mangaplus.shueisha.co.jp" {
		return false
	}
	parts := strings.Split(u.Path, "/")[1:]
	if parts[0] != "titles" {
		return false
	}
	s.MangaPlusID = parts[1]
	return true
}
