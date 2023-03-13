package penyusutan

import "time"

type UseCase interface {
	CalcPenyusutan(int, string, time.Time) error
}
