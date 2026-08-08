package seriesupdate

import (
	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/services/datasource"
	"github.com/abibby/icbmdb/services/datasource/mangadex"
	"github.com/davecgh/go-spew/spew"
)

func Update() error {
	datasource.RegisterDatasource(mangadex.New())
	// datasource.RegisterDatasource(mangaplus.New())

	s := &models.Series{
		DatasourceIDs: models.DatasourceIDs{
			MangadexID: "68112dc1-2b80-4f20-beb8-2f2a8716a430",
		},
	}
	i := &models.Issue{}
	for _, d := range datasource.Datasources() {
		id := d.ID(&s.DatasourceIDs)
		if id == "" {
			continue
		}

		applier, err := d.Series(id)
		if err != nil {
			return err
		}

		err = applier.ApplyData(s)
		if err != nil {
			return err
		}

		issuesAppliers, err := applier.Issues()
		if err != nil {
			return err
		}

		for _, issueApplier := range issuesAppliers[:1] {
			err = issueApplier.ApplyData(i)
			if err != nil {
				return err
			}
		}

		d.UpdateTimestamp(&s.DatasourceIDs)
	}

	spew.Dump(s)
	spew.Dump(i)
	return nil
}
