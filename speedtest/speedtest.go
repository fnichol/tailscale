package speedtest

import "time"

const (
	Start = "start" // Start the test.
	End   = "end"   // End the test.
	Data  = "data"  // Message contains data.

	LenBufJSON = 100   // agreed upon before hand. Buffer size for json messages.
	LenBufData = 32000 // buffer size for random bytes `
)

var (
	downloadTestDuration time.Duration = time.Second * 5
)

type Header struct {
	Type         string `json: "type"`
	IncomingSize int    `json: "incoming_size,omitempty"` // Currently being set but not read
}

type Record struct {
	TimeSlot time.Duration
	Size     int32
}
