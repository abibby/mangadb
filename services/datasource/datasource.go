package datasource

import (
	"slices"
	"sync"

	"github.com/abibby/icbmdb/app/models"
)

type Datasource interface {
	// Quality 100s community 200s distributor 300s publisher
	Quality() int

	Series(id string) (SeriesApplier, error)
	ApplySeriesID(s *models.Series, uri string) bool

	ID(s *models.DatasourceIDs) string
	UpdateTimestamp(s *models.DatasourceIDs)

	// ApplyIssueData(id string, s *models.Issue) error
	// IssueID(s *models.Issue) string
	// ApplyIssueID(s *models.Issue, uri string) bool
}

type Applier[T any] interface {
	ApplyData(s T) error
}
type SeriesApplier interface {
	Applier[*models.Series]

	Issues() ([]IssueApplier, error)
}
type IssueApplier interface {
	Applier[*models.Issue]
}

var sorted = true
var datasources = []Datasource{}
var mtx sync.RWMutex

func RegisterDatasource(d Datasource) {
	mtx.Lock()
	defer mtx.Unlock()
	sorted = false
	datasources = append(datasources, d)
}

func Datasources() []Datasource {
	if sorted == false {
		mtx.Lock()
		defer mtx.Unlock()
		slices.SortFunc(datasources, func(a, b Datasource) int {
			return a.Quality() - b.Quality()
		})
		sorted = true
	} else {
		mtx.RLock()
		defer mtx.RUnlock()
	}
	return datasources
}

func ApplySeriesID(s *models.Series, uri string) {
	for _, d := range Datasources() {
		d.ApplySeriesID(s, uri)
	}
}
