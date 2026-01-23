package common

import (
	"accounter/pkg/utils"
	"time"
)

// RequestParams request params
type RequestParams struct {
	DateStart int64   `db:"date_start" json:"date_start" form:"date_start"`
	DateEnd   int64   `db:"date_end" json:"date_end" form:"date_end"`
	Timezone  string  `db:"-" json:"timezone" form:"timezone"`
	Users     []int64 `db:"users" json:"users" form:"users"`
	Status    string  `db:"status" json:"status" form:"status"`

	OrderBy   string `db:"order_by" json:"order_by" form:"order_by"`
	OrderDesc bool   `db:"order_desc" json:"order_desc" form:"order_desc"`
	Skip      int    `db:"skip" json:"skip" form:"skip"`
	Limit     int    `db:"limit" json:"limit" form:"limit"`
}

// NewRequestParams creates new RequestParams
func NewRequestParams(n time.Time) RequestParams {
	return RequestParams{
		Timezone:  time.Local.String(),
		DateStart: time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()).Unix(),
		DateEnd:   time.Date(n.Year(), n.Month(), n.Day(), 23, 59, 59, 0, n.Location()).Unix(),
		OrderBy:   "date",
	}
}

// Encode encode params to tools.Params
func (tp RequestParams) Encode() utils.Params {
	params := utils.NewParams()

	params.Set("date_start", tp.DateStart)
	params.Set("date_end", tp.DateEnd)
	params.Set("status", tp.Status)
	params.Set("timezone", tp.Timezone)

	for _, id := range tp.Users {
		params.Add("users", id)
	}

	return params
}
