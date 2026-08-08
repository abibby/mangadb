package mangadex

import (
	"fmt"
	"time"

	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/mangadexv5"
)

type Mangadex struct {
	client *mangadexv5.Client
}

var _ datasource.Datasource = (*Mangadex)(nil)

func New() *Mangadex {
	return &Mangadex{
		client: mangadexv5.NewClient(),
	}
}

// Quality implements [datasource.Datasource].
func (m *Mangadex) Quality() int {
	return 100
}

// ID implements [datasource.Datasource].
func (m *Mangadex) ID(s *models.DatasourceIDs) string {
	return s.MangadexID
}

// UpdateTimestamp implements [datasource.Datasource].
func (m *Mangadex) UpdateTimestamp(s *models.DatasourceIDs) {
	s.MangadexUpdated = time.Now()
}

// Series implements [datasource.Datasource].
func (m *Mangadex) Series(id string) (datasource.SeriesApplier, error) {
	seriesList, _, err := m.client.MangaList(&mangadexv5.MangaListRequest{
		IDs: []string{id},
	})
	if err != nil {
		return nil, err
	}
	if len(seriesList) == 0 {
		return nil, fmt.Errorf("no series")
	}
	series := seriesList[0]

	return &Series{
		client: m.client,
		series: series,
	}, nil
}

// ApplySeriesID implements [datasource.Datasource].
func (m *Mangadex) ApplySeriesID(s *models.Series, uri string) bool {
	return false
}
