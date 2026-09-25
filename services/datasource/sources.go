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
const MangaplusQuality = 51

const VizSource = "viz"
const VizQuality = 50

type Series struct {
	Source string
	ID     string
}

func Get(u string) (*Series, bool) {
	if id, ok := strings.CutPrefix(u, "https://mangaplus.shueisha.co.jp/titles/"); ok {
		return &Series{
			Source: MangaplusSource,
			ID:     id,
		}, true
	}
	if id, ok := strings.CutPrefix(u, "https://www.viz.com/shonenjump/chapters/"); ok {
		return &Series{
			Source: VizSource,
			ID:     id,
		}, true
	}
	if rest, ok := strings.CutPrefix(u, "https://anilist.co/manga/"); ok {
		parts := strings.SplitN(rest, "/", 2)
		return &Series{
			Source: AnilistSource,
			ID:     parts[0],
		}, true
	}
	if rest, ok := strings.CutPrefix(u, "https://mangadex.org/title/"); ok {
		parts := strings.SplitN(rest, "/", 2)
		return &Series{
			Source: MangadexSource,
			ID:     parts[0],
		}, true
	}
	return nil, false
}
