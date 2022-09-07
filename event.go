package amplitude

// Event struct handler for event to be sent through Amplitude API
type Event struct {
	EventType          string         `json:"event_type"`
	UserID             string         `json:"user_id"`
	DeviceID           string         `json:"device_id"`
	Time               int64          `json:"time"`
	EventProperties    map[string]any `json:"event_properties"`
	UserProperties     map[string]any `json:"user_properties"`
	Groups             map[string]any `json:"groups"`
	AppVersion         string         `json:"app_version"`
	VersionName        string         `json:"version_name"`
	Library            string         `json:"library"`
	Platform           string         `json:"platform"`
	OsName             string         `json:"os_name"`
	OsVersion          string         `json:"os_version"`
	DeviceBrand        string         `json:"device_brand"`
	DeviceManufacturer string         `json:"device_manufacturer"`
	DeviceModel        string         `json:"device_model"`
	Carrier            string         `json:"carrier"`
	Country            string         `json:"country"`
	Region             string         `json:"region"`
	City               string         `json:"city"`
	Dma                string         `json:"dma"`
	Language           string         `json:"language"`
	UUID               string         `json:"uuid"`
	Price              float64        `json:"price"`
	Quantity           int            `json:"quantity"`
	Revenue            float64        `json:"revenue"`
	ProductID          string         `json:"productId"`
	RevenueType        string         `json:"revenueType"`
	LocationLat        float64        `json:"location_lat"`
	LocationLng        float64        `json:"location_lng"`
	IP                 string         `json:"ip"`
	Idfa               string         `json:"idfa"`
	Idfv               string         `json:"idfv"`
	Adid               string         `json:"adid"`
	AndroidID          string         `json:"android_id"`
	EventID            int            `json:"event_id"`
	SessionID          int64          `json:"session_id"`
	InsertID           string         `json:"insert_id"`
	Plan               *Plan          `json:"plan"`
}

// IsValid returns bool stating if Event is valid or not
func (e *Event) IsValid() bool {
	return e.DeviceID != "" || e.UserID != ""
}

// Plan struct handler for Plan premium feature
type Plan struct {
	Branch    string `json:"branch"`
	Source    string `json:"source"`
	Version   string `json:"version"`
	VersionID string `json:"versionId"`
}

type uploadEventInput struct {
	APIKey string   `json:"api_key"`
	Events []*Event `json:"events"`
}
