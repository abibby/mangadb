package datasource

import (
	"strings"
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

func AddFromURL(u string, ids map[string]string) {
	if id, ok := strings.CutPrefix(u, "https://mangaplus.shueisha.co.jp/titles/"); ok {
		ids[MangaplusSource] = id
	}
	if id, ok := strings.CutPrefix(u, "https://www.viz.com/shonenjump/chapters/"); ok {
		ids[VizSource] = id
	}
}
