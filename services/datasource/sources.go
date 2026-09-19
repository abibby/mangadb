package datasource

import (
	"strings"

	"abibby.com/mangadb/app/events"
)

// Quality
// scanlation 0-24
// fan db 25-49
// official 50-75

const MangadexSource = "mangadex"
const MangadexQuality = 1

const AnilistSource = "anilist"
const AnilistQuality = 25

const MangaplusSource = "mangaplus"
const MangaplusQuality = 50

const VizSource = "viz"
const VizQuality = 51

func Get(u string) (*events.FetchSeriesEvent, bool) {
	if id, ok := strings.CutPrefix(u, "https://mangaplus.shueisha.co.jp/titles/"); ok {
		return &events.FetchSeriesEvent{
			Source: MangaplusSource,
			ID:     id,
		}, true
	}
	if id, ok := strings.CutPrefix(u, "https://www.viz.com/shonenjump/chapters/"); ok {
		return &events.FetchSeriesEvent{
			Source: VizSource,
			ID:     id,
		}, true
	}
	if rest, ok := strings.CutPrefix(u, "https://anilist.co/manga/"); ok {
		parts := strings.SplitN(rest, "/", 2)
		return &events.FetchSeriesEvent{
			Source: AnilistSource,
			ID:     parts[0],
		}, true
	}
	return nil, false
}
