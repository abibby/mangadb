package mangaplus

import (
	"strconv"
	"strings"
	"time"

	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/icbmdb/services/mangaplus/mpproto"
)

type Issue struct {
	chapter *mpproto.Chapter
}

var _ datasource.Applier[*models.Issue] = (*Issue)(nil)

// ApplyData implements [datasource.Applier].
func (a *Issue) ApplyData(i *models.Issue) error {

	num, err := strconv.ParseFloat(strings.TrimPrefix(a.chapter.Name, "#"), 32)
	if err != nil {
		return err
	}

	i.Number = float32(num)
	i.Title = strings.SplitAfterN(a.chapter.GetSubTitle(), ": ", 2)[1]
	i.ReleaseDate = time.Unix(a.chapter.GetStartTimeStamp(), 0)
	return nil
}

// (*mpproto.Chapter)(0xc000355cb0)(
// 	titleId:100171
// 	chapterId:1009921
// 	name:"#001"
// 	subTitle:"Chapter 1: That's How Love Starts, Ya Know!"
// 	thumbnailUrl:"https://jumpg-assets.tokyo-cdn.com/secure/title/100171/chapter/1009921/chapter_thumbnail/178402.jpg?hash=CX8L1Z4jWNWkcazZ0yE4Tw&expires=1786186800"
// 	startTimeStamp:1629730800
// 	endTimeStamp:2145884400
// 	viewCount:1967071
// 	commentCount:918
// )
