package layer

type FetchActionBase interface {
	GetGlobalTimeline() []Note
	GetLocalTimeline() []Note
	GetHomeTimeline() []Note

	GetPost(id string) Note

	GetNotifications() []Notification
	GetNotification(id string) Notification

	GetUser(id string) User
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

func (f *DataFetch) GetPost(id string) Note {
	return f.base.GetPost(id)
}

func (f *DataFetch) GetUser(id string) User {
	return f.base.GetUser(id)
}

func (f *DataFetch) GetNotification(id string) Notification {
	return f.base.GetNotification(id)
}

func (f *DataFetch) GetNotifications(id string) []Notification {
	return f.base.GetNotifications()
}

type SendActionBase interface {
}

type DataSend struct {
	base FetchActionBase
}

type DeleteActionBase interface {
}

type OtherActionBase interface {
}
