package models

import "time"

type Banner struct {
	Tags       []int     `json:"tag"`
	Fid        int       `json:"fid"`
	Bid        int       `json:"bid"`
	Cnt        Content   `json:"content,omitempty"`
	Is_active  bool      `json:"is_active"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Content struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Pair struct {
	Tag int
	Fid int
}

type DbBanner struct {
	Tag        int
	Fid        int
	Bid        int
	Cnt        Content
	Is_active  bool
	Created_at time.Time
	Updated_at time.Time
}
