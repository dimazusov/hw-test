package storage

// GettingEventParams structure for getting events
type GettingEventParams struct {
	IDs         []uint
	Title       string
	FromTime    uint
	ToTime      uint
	ExactTime   uint
	Timezone    uint
	UserID      uint
	Page        uint
	CountOnPage uint
}
