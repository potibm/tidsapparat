package formatters

type JSGroup struct {
	Type    string    `json:"@type"` // Has to be "Group"
	UID     string    `json:"uid"`
	ProdID  string    `json:"prodId,omitempty"`
	Updated string    `json:"updated"`
	Entries []JSEvent `json:"entries"`
}

type JSEvent struct {
	Type        string                `json:"@type"` // Has to be "jsevent"
	UID         string                `json:"uid"`
	ProdID      string                `json:"prodId,omitempty"`
	Created     string                `json:"created,omitempty"` // UTC
	Updated     string                `json:"updated,omitempty"` // UTC
	Sequence    uint64                `json:"sequence,omitempty"`
	Title       string                `json:"title"`
	Description string                `json:"description,omitempty"`
	Start       string                `json:"start"`
	TimeZone    string                `json:"timeZone,omitempty"`
	Duration    string                `json:"duration,omitempty"` // ISO 8601 (e.g. "PT1H30M")
	Status      string                `json:"status,omitempty"`
	Locations   map[string]JSLocation `json:"locations,omitempty"`
	Links       map[string]JSLink     `json:"links,omitempty"`
	Keywords    map[string]bool       `json:"keywords,omitempty"`
}

type JSLocation struct {
	Type        string `json:"@type"` // "Location"
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type JSLink struct {
	Type string `json:"@type"` // "Link"
	Href string `json:"href"`
}
