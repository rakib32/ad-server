package responses

type Ad struct {
	Adid          int64
	AdName        string
	AppStoreURL   string
	ClickURL      string
	ImpressionURL string
	ActionText    string
	Images        Image
}
