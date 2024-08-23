package xsyslog

// RFC3339 Time Formatter. Support NanoUTC and MicroUTC.
const (
	RFC3339NanoUTC  = "2006-01-02T15:04:05.999999999Z"
	RFC3339MicroUTC = "2006-01-02T15:04:05.999999Z"
)

type Event struct {
	Priority  uint8  `json:"priority"`
	ProcID    string `json:"proc_id"`
	Level     string `json:"level"`
	Tag       string `json:"tag"`
	AppName   string `json:"stream"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	//Tags      map[string]string `json:"tags"`
}
