package models

import (
	"context"
	"time"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
	"github.com/google/uuid"
)

//go:generate spice generate:migration
type Issue struct {
	BaseModel

	// Title of the book.
	Title string `json:"title" db:"title"`

	// Title of the series the book is part of.
	SeriesID uuid.UUID `json:"series_id" db:"series_id"`

	// Number of the book in the series.
	Number float32 `json:"number" db:"number"`

	// Volume containing the book.
	//
	// Volume is a notion that is specific to US
	// Comics, where the same series can have multiple volumes. Volumes can be
	// referenced by number (1, 2, 3…) or by year (2018, 2020…).
	Volume int `json:"volume" db:"volume"`

	// VolumeID uuid.UUID `json:"volume_id" db:"volume_id"`

	// Quite specific to US comics, some books can be part of cross-over sory
	// arcs. Those fields can be used to specify an alternate series, its number
	// and count of books.
	//
	// AlternateSeries / AlternateNumber / AlternateCount

	// A description or summary of the book.
	Summary string `json:"summary" db:"summary"`

	// A free text field, usually used to store information about the
	// application that created the ComicInfo.xml file.
	Notes string `json:"notes" db:"notes"`

	// Year / Month / Day
	// Usually contains the release date of the book.
	ReleaseDate time.Time `json:"release_date" db:"release_date"`

	// An imprint is a group of publications under the umbrella of a larger
	// imprint or a Publisher. For example, Vertigo is an Imprint of DC Comics.
	Imprint string `json:"imprint" db:"imprint"`

	// A URL pointing to a reference website for the book.
	//
	// It is accepted that multiple values are space separated. If a space is a
	// part of the url it must be percent encoded.
	Web string `json:"web" db:"web"`

	// A language code describing the language of the book.
	//
	// Without any information on what kind of code this element is supposed to
	// contain, it is recommended to use the IETF BCP 47 language tag, which can
	// describe the language but also the script used. This helps to
	// differentiate languages with multiple scripts, like Traditional and
	// Simplified Chinese.
	//
	// See also:
	//
	//     Choosing a language tag - W3C
	//     Language subtag lookup app
	LanguageISO string `json:"language_iso" db:"language_iso"`

	// The original publication's binding format for scanned physical books or
	// presentation format for digital sources.
	//
	// "TBP", "HC", "Web", "Digital" are common designators.
	Format string `json:"format" db:"format"`

	// Series builder.BelongsTo[*Series] `json:"series" db:"series"`
	// Volume builder.BelongsTo[*Volume] `json:"volume" db:"volume"`
	// Staff  builder.HasMany[*Staff]    `json:"staff " db:"Staff "`

	DatasourceIDs
}

func init() {
	providers.Add(modeldi.Register[*Issue])
}

func IssueQuery(ctx context.Context) *builder.ModelBuilder[*Issue] {
	return builder.From[*Issue]().WithContext(ctx)
}
