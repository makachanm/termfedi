package layer

type FetchActionBase interface {
	GetGlobalTimeline() []Note
	GetLocalTimeline() []Note
	GetHomeTimeline() []Note

	GetNotifications() []Notification

	PostRenote(note_id string) bool
	PostReaction(note_id string) bool
}

type DataFetch struct {
	base FetchActionBase
}

func NewDataFetchAction(b FetchActionBase) DataFetch {
	return DataFetch{base: b}
}

func (f *DataFetch) GetGlobalTimeline() []Note {
	return f.base.GetGlobalTimeline()
}

func (f *DataFetch) GetLocalTimeline() []Note {
	return f.base.GetLocalTimeline()
}

func (f *DataFetch) GetHomeTimeline() []Note {
	return f.base.GetHomeTimeline()
}

func (f *DataFetch) GetNotifications() []Notification {
	return f.base.GetNotifications()
}

func (f *DataFetch) PostRenote(note_id string) bool {
	return f.base.PostRenote(note_id)
}

func (f *DataFetch) PostReaction(note_id string) bool {
	return f.base.PostReaction(note_id)
}
