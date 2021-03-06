package domain

type Event struct {
	ID               uint   `json:"id" db:"id"`
	Title            string `json:"title" db:"title"`
	Time             uint   `json:"time" db:"time"`
	Timezone         uint8  `json:"timezone" db:"timezone"`
	Duration         uint   `json:"duration" db:"duration"`
	Description      string `json:"description" db:"description"`
	UserID           uint   `json:"userId" db:"user_id"`
	NotificationTime uint   `json:"notificationTime" db:"notification_time"`
}
