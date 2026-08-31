package datasource

import (
	"fmt"
	"sync"

	"github.com/abibby/icbmdb/app/models"
)

type Datasource interface {
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

	Key() ChapterKey
}

type ChapterKey string

func NewChapterKey(series string, chapter float32) ChapterKey {
	return ChapterKey(fmt.Sprintf("%s-%d", series, chapter))
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
	return datasources
}
