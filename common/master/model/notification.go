package model

type Notification struct {
	Text      string `json:"text"`
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	SoundName string `json:"sound_name,omitempty"`
}
