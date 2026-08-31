package mangadex

import (
	"strconv"

	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/mangadexv5"
)

type Issue struct {
	chapter *mangadexv5.Chapter
}

var _ datasource.IssueApplier = (*Issue)(nil)

// ApplyData implements [datasource.Applier].
func (a *Issue) ApplyData(i *models.Issue) error {
	i.Title = a.chapter.Title

	if num, ok := a.number(); ok {
		i.Number = num
	}

	if vol, err := strconv.Atoi(a.chapter.Volume.String()); err == nil {
		i.Volume = vol
	}
	i.ReleaseDate = a.chapter.PublishAt

	return nil
}

// Key implements [datasource.IssueApplier].
func (a *Issue) Key() datasource.ChapterKey {
	num, _ := a.number()
	return datasource.NewChapterKey("", num)
}

func (a *Issue) number() (float32, bool) {
	num, err := strconv.ParseFloat(a.chapter.Chapter.String(), 32)
	if err != nil {
		return 0, false
	}
	return float32(num), true
}

// (*mangadexv5.Chapter)(0xc00045a000)({
//  Model: (mangadexv5.Model) {
//   ID: (string) (len=36) "4c03a97e-0688-4fef-93d8-df942b777acc",
//   Relationships: (mangadexv5.RelationshipList) (len=3 cap=4) {
//    (*mangadexv5.Relationship)(0xc0002109c0)({
//     ID: (string) (len=36) "4f1de6a2-f0c5-4ac5-bce5-02c7dbb67deb",
//     Type: (string) (len=16) "scanlation_group"
//    }),
//    (*mangadexv5.Relationship)(0xc0002109e0)({
//     ID: (string) (len=36) "68112dc1-2b80-4f20-beb8-2f2a8716a430",
//     Type: (string) (len=5) "manga"
//    }),
//    (*mangadexv5.Relationship)(0xc000210a20)({
//     ID: (string) (len=36) "74d95af1-7492-4fca-bc44-10c9142703e8",
//     Type: (string) (len=4) "user"
//    })
//   }
//  },
//  Title: (string) (len=32) "That's How Love Starts, Ya Know!",
//  Volume: (*nulls.String)(0xc000536020)((len=1) 1),
//  Chapter: (*nulls.String)(0xc000536030)((len=1) 1),
//  Pages: (int) 0,
//  TranslatedLanguage: (string) (len=2) "en",
//  Uploader: (string) "",
//  ExternalURL: (*nulls.String)(0xc000536040)((len=47) https://mangaplus.shueisha.co.jp/viewer/1009921),
//  Version: (int) 4,
//  CreatedAt: (time.Time) 2021-11-09 01:54:47 +0000 +0000,
//  UpdatedAt: (time.Time) 2022-08-29 16:04:55 +0000 +0000,
//  PublishAt: (time.Time) 2021-11-09 01:54:47 +0000 +0000,
//  ReadableAt: (time.Time) 2021-11-09 01:54:47 +0000 +0000,
//  manga: (*mangadexv5.Manga)(<nil>)
// })
