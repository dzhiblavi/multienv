package model

// osascript -e 'display notification "The backup script finished successfully." with title "Backup Status" subtitle "Success" sound name "Glass"'
type Notification struct {
	Text      string `json:"text"`
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	SoundName string `json:"sound_name,omitempty"`
}
