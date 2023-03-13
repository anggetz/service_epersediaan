package pemeliharaan

import (
	"pvg/simada/service-golang/domains"
	"time"
)

type Model struct {
	tableName struct{}   `pg:"pemeliharaan,discard_unknown_columns"`
	ID        int        `json:"id"`
	Biaya     float64    `json:"biaya"`
	Tgl       *time.Time `json:"tgl"`
	domains.GenericModel
}
