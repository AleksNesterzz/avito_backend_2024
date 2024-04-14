package csv_log

import "time"

type CsvLog struct {
	Client    int
	Operation string
	Segment   string
	Dt        time.Time
}
