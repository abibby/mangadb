package mangadex

import (
	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/mangadexv5"
	"github.com/abibby/salusa/slices"
)

type Series struct {
	client *mangadexv5.Client
	series *mangadexv5.Manga
}

var _ datasource.SeriesApplier = (*Series)(nil)

// ApplyData implements [datasource.SeriesApplier].
func (a *Series) ApplyData(s *models.Series) error {

	s.Title = a.series.Title.String()
	s.Aliases = slices.Map(a.series.AltTitles, mangadexv5.LangMap.String)

	datasource.ApplySeriesID(s, a.series.Links.EnglishTranslationURL)
	return nil
}

// Issues implements [datasource.SeriesApplier].
func (s *Series) Issues() ([]datasource.IssueApplier, error) {

	chapters, _, err := s.client.ChapterList(&mangadexv5.ChapterListRequest{
		MangaID: s.series.ID,
	})
	if err != nil {
		return nil, err
	}

	issues := make([]datasource.IssueApplier, len(chapters))
	for i, c := range chapters {
		issues[i] = &Issue{
			chapter: c,
		}
	}
	return issues, nil
}
