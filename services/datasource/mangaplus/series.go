package mangaplus

// type Series struct {
// 	client      *mangaplus.Client
// 	titleDetail *mpproto.TitleDetailView
// }

// var _ datasource.SeriesApplier = (*Series)(nil)

// // ApplyData implements [datasource.SeriesApplier].
// func (a *Series) ApplyData(s *models.Series) error {
// 	title := a.titleDetail.GetTitle()
// 	s.Title = title.GetName()
// 	s.Description = a.titleDetail.GetOverview()
// 	return nil
// }

// // Issues implements [datasource.SeriesApplier].
// func (s *Series) Issues() ([]datasource.IssueApplier, error) {

// 	groups := s.titleDetail.GetChapterListGroup()

// 	issues := []datasource.IssueApplier{}
// 	for c := range chapters(groups) {
// 		issues = append(issues, &Issue{
// 			chapter: c,
// 		})
// 	}

// 	return issues, nil
// }

// func chapters(groups []*mpproto.TitleDetailView_ChapterGroup) iter.Seq[*mpproto.Chapter] {
// 	return func(yield func(*mpproto.Chapter) bool) {
// 		for _, group := range groups {
// 			for _, c := range group.FirstChapterList {
// 				if !yield(c) {
// 					return
// 				}
// 			}
// 			for _, c := range group.MidChapterList {
// 				if !yield(c) {
// 					return
// 				}
// 			}
// 			for _, c := range group.LastChapterList {
// 				if !yield(c) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }
