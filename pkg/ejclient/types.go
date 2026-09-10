package ejclient

import "strconv"

type RequestDateMode uint8

const (
	RequestDateModeDefault RequestDateMode = 0
	RequestDateModeISO     RequestDateMode = 1
	RequestDateModeUnix    RequestDateMode = 2
)

func (v RequestDateMode) String() string {
	return strconv.Itoa(int(v))
}

type Config struct {
	URL       string
	ContestID int32
	Token     string
}
