package penyusutan

import (
	"fmt"
	"log"
	"math"
	"pvg/simada/service-golang/domains"
	"pvg/simada/service-golang/domains/inventaris"
	"pvg/simada/service-golang/domains/kriteria_masa_manfaat"
	"pvg/simada/service-golang/domains/pemeliharaan"
	"strconv"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type UseCaseImpl struct{}

func NewUseCase() UseCase {
	return &UseCaseImpl{}
}

func (u *UseCaseImpl) CalcPenyusutan(opdId int, jenisAset string, dateRunning time.Time) error {
	// data, err := NewRepository[BasePenyusutan]().GetAll("report_penyusutan_2022", 1, 10, 0, "")

	// for _, d := range data {
	// 	fmt.Println("test id", d.ID)
	// }

	// if err != nil {
	// 	return err
	// }

	tanggalStartCalc := 2014
	tanggalEndCalc := dateRunning.Year()

	if dateRunning.Year() < 2014 {
		fmt.Errorf("Tidak boleh running lebih kecil dari 2014")
	}

	inventarises, err := domains.NewGenericRepository[*inventaris.Model]().All(-1, -1, func(q *orm.Query) {

		if opdId != 0 {
			q.Where("inventaris.pidopd = ?", opdId)
		}

		if jenisAset != "" {
			fmt.Println(jenisAset)
			q.Relation("Barang", func(q *orm.Query) (*orm.Query, error) {
				q.Where("barang.kode_jenis = ?", jenisAset)
				q.Where("barang.umur_ekonomis != ?", 0)

				return q, nil
			})
			// q.Where("barang.umur_ekonomis != ?", 0)

		} else {
			q.Relation("Barang")
		}
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(len(inventarises))

	mtx := sync.Mutex{}
	wg := sync.WaitGroup{}
	nextSleep := 50

	for index, inv := range inventarises {

		wg.Add(1)

		go func(inv *inventaris.Model, index int) {
			defer wg.Done()

			umurEkonomis := inv.Barang.UmurEkonomis

			if tanggalStartCalc < inv.TahunPerolehan {
				tanggalStartCalc = inv.TahunPerolehan
			}

			if tanggalEndCalc < inv.TahunPerolehan {
				return
			}
			for year := tanggalStartCalc; year <= tanggalEndCalc; year++ {
				func(umurEkonomis int, year int, inv *inventaris.Model, index int) {
					log.Println("running year", year, "index data -", index)

					needMigrate := index == 1

					bebanPenyusutan := 0.00
					monthRunning := 12
					prevYear := year - 1
					prevTableYear := "report_penyusutan_" + strconv.Itoa(prevYear)
					prevPenyusutan := BasePenyusutan{}

					today := time.Now()

					penyusutan := BasePenyusutan{
						InventarisID: inv.ID,
						CreatedDate:  &today,
						UpdatedDate:  &today,
					}
					pemeliharaanData := pemeliharaan.Model{}

					if inv.HargaSatuan != 0 {
						bebanPenyusutan = inv.HargaSatuan / float64(umurEkonomis)
					}

					if year == dateRunning.Year() {
						monthRunning = int(dateRunning.Month())
					}

					runningPenyusutan := time.Date(year, time.Month(monthRunning), 15, 0, 0, 0, 0, time.Local)
					penyusutan.RunningPenyusutan = &runningPenyusutan

					// check if data existing
					if ok, _ := domains.NewGenericRepository[*BasePenyusutan]().
						SetTableName("report_penyusutan_" + strconv.Itoa(year)).
						Migrate(needMigrate).
						Exists(func(q *orm.Query) {
							q.Where("inventaris_id = ?", inv.ID)
						}); ok {
						return
					}
					// get pemeliharaan
					pems, err := domains.NewGenericRepository[*pemeliharaan.Model]().All(1, 0, func(q *orm.Query) {
						q.ColumnExpr("sum(biaya) as biaya, MAX(tgl) tgl")
						q.Where("pidinventaris = ?", inv.ID)
						q.Where("is_exec_by_penyusutan IS FALSE ")
						q.Where("date_part('year', tgl) = ? ", year)
						q.OrderExpr("tgl DESC")
						q.GroupExpr("date_part('year', tgl)")
					})

					if err != nil && err != pg.ErrNoRows {
						log.Fatal("error when getting data pemeliharaan", err.Error())
						return
					}

					prevPenyusutans, err := domains.NewGenericRepository[*BasePenyusutan]().
						SetTableName(prevTableYear).
						Migrate(needMigrate).All(1, 0, func(q *orm.Query) {

						q.Where("inventaris_id = ?", inv.ID)
					})

					if err != nil && err != pg.ErrNoRows {
						log.Fatal("error when getting previous penyusutan", err.Error())
						return
					}

					if len(pems) > 0 {
						pemeliharaanData = *pems[0]
						penyusutan.TotalAtribusi = pemeliharaanData.Biaya
						penyusutan.BulanAtribusi = float64(pemeliharaanData.Tgl.Month())
					}

					if len(prevPenyusutans) > 0 {
						prevPenyusutan = *prevPenyusutans[0]
					}

					u.subCalcBebanPenyusutanSebelumAtribusi(&penyusutan, &prevPenyusutan, &pemeliharaanData, year)
					u.subCalcNilaiBukuSebelumAtribusi(&penyusutan, &prevPenyusutan, year)
					u.subCalcNilaiBukuSetelahAtribusi(&penyusutan, &prevPenyusutan, &pemeliharaanData, year)
					u.subCalcPenambahAtribusi(&penyusutan, year)
					u.subCalcPemakaianSDTahunBerkenaan(&penyusutan, inv)
					u.subCalcSisaUmurEkonomisSDTahunSebelumnya(&penyusutan, &prevPenyusutan, inv, year)
					u.subCalcSisaUmurEkonomisSDBulanBerkenaan(&penyusutan, year)
					u.subCalcPenambahUmurEkonomis(&penyusutan, inv, year)
					u.subSisaUmurEkonomisSetelahAtribusi(&penyusutan, umurEkonomis, year)
					u.subCalcBebanPenyusutanSetelahAtribusi(&penyusutan, &prevPenyusutan, inv, umurEkonomis, bebanPenyusutan, year)
					u.subCalcBebanPenyusutanTahunBerkenaan(&penyusutan, inv, year)
					u.subCalcAkumulasiTahunBerkenaan(&penyusutan, &prevPenyusutan, inv, bebanPenyusutan, year)
					u.subCalcNilaiBuku(&penyusutan, &prevPenyusutan, inv, year)
					mtx.Lock()
					err = domains.NewGenericRepository[*BasePenyusutan]().SetTableName("report_penyusutan_" + strconv.Itoa(year)).Insert(&penyusutan)

					defer mtx.Unlock()
					if err != nil {
						log.Println("Error when insert data", err)
						return
					}

					log.Println("inserted the data")

				}(umurEkonomis, year, inv, index)
			}
		}(inv, index)

		if index == nextSleep {
			wg.Wait()

			log.Println("sleep for awhile")

			nextSleep += 50
		}
	}

	return nil

}

func (u *UseCaseImpl) subCalcNilaiBuku(peny *BasePenyusutan, prevPeny *BasePenyusutan, inv *inventaris.Model, year int) {
	if year >= 2015 {
		if inv.TahunPerolehan > year {
			peny.NilaiBukuTahunBerkenaan = 0
		} else {
			if inv.TahunPerolehan == year {
				peny.NilaiBukuTahunBerkenaan = inv.HargaSatuan - peny.BebanPenyusutanTahunBerkenaan
			} else {
				peny.NilaiBukuTahunBerkenaan = (prevPeny.NilaiBukuTahunBerkenaan + peny.TotalAtribusi) - peny.BebanPenyusutanTahunBerkenaan
			}
		}
	} else {
		if inv.TahunPerolehan > year {
			peny.NilaiBukuTahunBerkenaan = 0
		} else {
			peny.NilaiBukuTahunBerkenaan = (inv.HargaSatuan + peny.TotalAtribusi) - peny.AkumulasiPenyusutanSDTahunBerkenaan
		}
	}

	if peny.NilaiBukuTahunBerkenaan < 0 {
		peny.NilaiBukuTahunBerkenaan = 0
	}
}

func (u *UseCaseImpl) subCalcAkumulasiTahunBerkenaan(peny *BasePenyusutan, prevPreny *BasePenyusutan, inv *inventaris.Model, bebanPenyusutan float64, year int) {
	if year < 2015 {

		if inv.TahunPerolehan == 2014 {
			peny.AkumulasiPenyusutanSDTahunBerkenaan = ((float64(peny.RunningPenyusutan.Month()+1) - float64(inv.TglPerolehan.Month())) * bebanPenyusutan)
		} else {

			akumulasi1 := float64(peny.PemakaianSDTahunBerkenaan) * peny.BebanPenyusutanSetelahAtribusi
			akumulasi2 := inv.HargaSatuan + peny.TotalAtribusi

			log.Println("here brooss? #2", akumulasi1, akumulasi2)

			if akumulasi1 > akumulasi2 {
				peny.AkumulasiPenyusutanSDTahunBerkenaan = akumulasi2
			} else {
				peny.AkumulasiPenyusutanSDTahunBerkenaan = akumulasi1
			}
		}
	} else {
		peny.AkumulasiPenyusutanSDTahunBerkenaan = peny.BebanPenyusutanTahunBerkenaan + prevPreny.AkumulasiPenyusutanSDTahunBerkenaan
	}

	if peny.AkumulasiPenyusutanSDTahunBerkenaan < 0 {
		peny.AkumulasiPenyusutanSDTahunBerkenaan = 0
	}

	log.Println("data for insert inside", peny.AkumulasiPenyusutanSDTahunBerkenaan, year)
}

func (u *UseCaseImpl) subCalcBebanPenyusutanTahunBerkenaan(peny *BasePenyusutan, inv *inventaris.Model, year int) {
	if year >= 2015 {
		if inv.TahunPerolehan > year {
			peny.BebanPenyusutanTahunBerkenaan = 0
		} else {
			if inv.TahunPerolehan == year && int(peny.RunningPenyusutan.Month()) < int(inv.TglPerolehan.Month()) {
				peny.BebanPenyusutanTahunBerkenaan = 0
			} else {
				if inv.TahunPerolehan == year && int(peny.RunningPenyusutan.Month()) >= int(inv.TglPerolehan.Month()) {
					peny.BebanPenyusutanTahunBerkenaan = (((float64(peny.RunningPenyusutan.Month()) + 1) - float64(inv.TglPerolehan.Month())) * peny.BebanPenyusutanSetelahAtribusi)
				} else {
					if peny.TotalAtribusi != 0 && int(peny.RunningPenyusutan.Month()) >= int(peny.BulanAtribusi) {
						peny.BebanPenyusutanTahunBerkenaan = (peny.BebanPenyusutanSebelumAtribusi + (((float64(peny.RunningPenyusutan.Month()) + 1) - float64(inv.TglPerolehan.Month())) * peny.BebanPenyusutanSetelahAtribusi))

					} else if peny.NilaiBukuSetelahBulanAtribusi < float64(peny.RunningPenyusutan.Month())*peny.BebanPenyusutanSetelahAtribusi {
						peny.BebanPenyusutanTahunBerkenaan = peny.BebanPenyusutanSetelahAtribusi
					} else {
						peny.BebanPenyusutanTahunBerkenaan = float64(peny.RunningPenyusutan.Month()) * peny.BebanPenyusutanSetelahAtribusi
					}
				}
			}

		}
	}

	peny.BebanPenyusutanTahunBerkenaan = math.Ceil(peny.BebanPenyusutanTahunBerkenaan)

	if peny.BebanPenyusutanTahunBerkenaan < 0 {
		peny.BebanPenyusutanTahunBerkenaan = 0
	}
}

func (u *UseCaseImpl) subCalcBebanPenyusutanSetelahAtribusi(peny *BasePenyusutan, prevPeny *BasePenyusutan, inv *inventaris.Model, umurEkonomis int, bebanPenyusutan float64, year int) {
	if year >= 2015 {
		if inv.TahunPerolehan > year {
			peny.BebanPenyusutanSetelahAtribusi = 0
		} else {
			if inv.TahunPerolehan == year {
				peny.BebanPenyusutanSetelahAtribusi = bebanPenyusutan
			} else {
				if int(peny.RunningPenyusutan.Month()) < int(peny.BulanAtribusi) {
					peny.BebanPenyusutanSetelahAtribusi = prevPeny.BebanPenyusutanSetelahAtribusi
				} else if peny.TotalAtribusi != 0 {
					if peny.NilaiBukuSetelahBulanAtribusi == 0 || peny.SisaUmurEkonomisSetelahAtribusi == 0 {
						peny.BebanPenyusutanSetelahAtribusi = 0
					} else {
						peny.BebanPenyusutanSetelahAtribusi = peny.NilaiBukuSetelahBulanAtribusi / float64(peny.SisaUmurEkonomisSetelahAtribusi)
					}
				} else {
					peny.BebanPenyusutanSetelahAtribusi = prevPeny.BebanPenyusutanSetelahAtribusi
				}
			}
		}
	} else if year == 2014 {
		peny.BebanPenyusutanSetelahAtribusi = (inv.HargaSatuan + peny.TotalAtribusi) / float64(umurEkonomis)
	}

	peny.BebanPenyusutanSetelahAtribusi = math.Ceil(peny.BebanPenyusutanSetelahAtribusi)

	if peny.BebanPenyusutanSetelahAtribusi < 0 {
		peny.BebanPenyusutanSetelahAtribusi = 0
	}
}

func (u *UseCaseImpl) subSisaUmurEkonomisSetelahAtribusi(peny *BasePenyusutan, umurEkonomis int, year int) {
	if year >= 2015 {
		if peny.TotalAtribusi == 0 {
			peny.SisaUmurEkonomisSetelahAtribusi = peny.SisaUmurEkonomisSDBulanBerkenaan
		} else {
			if peny.SisaUmurEkonomisSDBulanBerkenaan == 0 {
				peny.SisaUmurEkonomisSetelahAtribusi = peny.SisaUmurEkonomisSDBulanBerkenaan + peny.PenambahanUmurEkonomis
			} else {
				akumulasiFormula := (peny.SisaUmurEkonomisSDBulanBerkenaan - (int(peny.BulanAtribusi) - 1)) + peny.PenambahanUmurEkonomis
				if akumulasiFormula > umurEkonomis {
					akumulasiFormula = umurEkonomis
				}

				peny.SisaUmurEkonomisSetelahAtribusi = akumulasiFormula
			}
		}

		if peny.SisaUmurEkonomisSetelahAtribusi > umurEkonomis {
			peny.SisaUmurEkonomisSetelahAtribusi = umurEkonomis
		}
	}
}

func (u *UseCaseImpl) subCalcPenambahUmurEkonomis(peny *BasePenyusutan, inv *inventaris.Model, year int) error {
	if year >= 2015 {
		tambahanMasaManfaat, err := domains.NewGenericRepository[*kriteria_masa_manfaat.Model]().All(1, 0, func(q *orm.Query) {
			q.Where(" ? LIKE '%' || kode_barang || '%'", inv.KodeBarang)
			q.Where(" ? between min and maks", peny.PenambahanAtribusi)
		})

		if err != nil || err == pg.ErrNoRows {
			return err
		}

		if len(tambahanMasaManfaat) > 0 {
			peny.PenambahanUmurEkonomis = tambahanMasaManfaat[0].TahunTambahan * 12
		}

	}

	return nil
}

func (u *UseCaseImpl) subCalcSisaUmurEkonomisSDBulanBerkenaan(peny *BasePenyusutan, year int) {
	if year >= 2015 {
		if int(peny.RunningPenyusutan.Month()) > int(peny.BulanAtribusi-1) {
			if peny.TotalAtribusi == 0 {
				peny.SisaUmurEkonomisSDBulanBerkenaan = peny.SisaUmurEkonomisSDTahunSebelumnya
			} else {
				akumulasiFormula := peny.SisaUmurEkonomisSDTahunSebelumnya - (int(peny.BulanAtribusi) - 1)
				if akumulasiFormula < 0 {
					akumulasiFormula = peny.SisaUmurEkonomisSDTahunSebelumnya
				}

				peny.SisaUmurEkonomisSDBulanBerkenaan = akumulasiFormula
			}
		} else {
			peny.SisaUmurEkonomisSDBulanBerkenaan = peny.SisaUmurEkonomisSDTahunSebelumnya
		}
	}
}

func (u *UseCaseImpl) subCalcSisaUmurEkonomisSDTahunSebelumnya(peny *BasePenyusutan, prevPeny *BasePenyusutan, inv *inventaris.Model, year int) {
	if year >= 2015 {
		if inv.TahunPerolehan > year-1 {
			peny.SisaUmurEkonomisSDTahunSebelumnya = 0
		} else {
			if year == 2015 {
				peny.SisaUmurEkonomisSDTahunSebelumnya = inv.Barang.UmurEkonomis - prevPeny.PemakaianSDTahunBerkenaan
			} else {
				if inv.TahunPerolehan == year {
					peny.SisaUmurEkonomisSDTahunSebelumnya = inv.Barang.UmurEkonomis - 13 - int(inv.TglPerolehan.Month())
				} else {
					substract := 12
					if prevPeny.SisaUmurEkonomisSetelahAtribusi < substract {
						peny.SisaUmurEkonomisSDTahunSebelumnya = 0
					} else {
						akumulasiFormula := 12
						if prevPeny.BulanAtribusi != 0 {
							akumulasiFormula = 13 - int(prevPeny.BulanAtribusi)
						}

						peny.SisaUmurEkonomisSDTahunSebelumnya = prevPeny.SisaUmurEkonomisSetelahAtribusi - akumulasiFormula
					}
				}
			}
		}
	}
}

func (u *UseCaseImpl) subCalcPenambahAtribusi(peny *BasePenyusutan, year int) {

	if year >= 2015 {
		if peny.TotalAtribusi == 0 {
			peny.PenambahanAtribusi = 0
		} else {
			if peny.NilaiBukuSebelumBulanAtribusi == 0 {
				peny.PenambahanAtribusi = 0
			} else {
				peny.PenambahanAtribusi = (peny.TotalAtribusi / peny.NilaiBukuSebelumBulanAtribusi) * 100
			}
		}
	}

	if peny.PenambahanAtribusi > 100 {
		peny.PenambahanAtribusi = 100
	}
}

func (u *UseCaseImpl) subCalcPemakaianSDTahunBerkenaan(peny *BasePenyusutan, inv *inventaris.Model) {
	durationDiff := peny.RunningPenyusutan.Sub(*inv.TglPerolehan)
	monthDuration := math.Round(durationDiff.Hours() / 24 / 30)

	if peny.RunningPenyusutan.After(*inv.TglPerolehan) && monthDuration >= float64(inv.Barang.UmurEkonomis) {
		peny.PemakaianSDTahunBerkenaan = inv.Barang.UmurEkonomis
	} else {
		peny.PemakaianSDTahunBerkenaan = int(monthDuration)
	}
}

func (u *UseCaseImpl) subCalcNilaiBukuSetelahAtribusi(peny *BasePenyusutan, prevPeny *BasePenyusutan, pem *pemeliharaan.Model, year int) {
	if year >= 2015 {
		if int(peny.RunningPenyusutan.Month()) > int(peny.BulanAtribusi-1) && peny.BulanAtribusi != 0 {
			peny.NilaiBukuSetelahBulanAtribusi = peny.TotalAtribusi + peny.NilaiBukuSebelumBulanAtribusi
			log.Println(peny.NilaiBukuSetelahBulanAtribusi, year, "setelah bulan atribusi")
		} else {
			peny.NilaiBukuSetelahBulanAtribusi = peny.NilaiBukuSebelumBulanAtribusi
		}

		peny.NilaiBukuSetelahBulanAtribusi = math.Ceil(peny.NilaiBukuSetelahBulanAtribusi)
	}

	if peny.NilaiBukuSetelahBulanAtribusi < 0 {
		peny.NilaiBukuSetelahBulanAtribusi = 0
	}
}

func (u *UseCaseImpl) subCalcNilaiBukuSebelumAtribusi(peny *BasePenyusutan, prevPeny *BasePenyusutan, year int) {

	if year >= 2015 {
		if prevPeny.NilaiBukuTahunBerkenaan != 0 {
			peny.NilaiBukuSebelumBulanAtribusi = prevPeny.NilaiBukuTahunBerkenaan
		} else {
			peny.NilaiBukuSebelumBulanAtribusi = peny.BebanPenyusutanSebelumAtribusi
		}

		if peny.NilaiBukuSebelumBulanAtribusi < 0 {
			peny.NilaiBukuSebelumBulanAtribusi = 0
		}
	}

}

func (u *UseCaseImpl) subCalcBebanPenyusutanSebelumAtribusi(peny *BasePenyusutan, prevPeny *BasePenyusutan, pem *pemeliharaan.Model, year int) {

	if year >= 2015 {
		if peny.TotalAtribusi != 0 {
			if int(peny.RunningPenyusutan.Month()) > int(peny.BulanAtribusi) {
				if peny.TotalAtribusi == 0 {
					peny.BebanPenyusutanSebelumAtribusi = 0
				} else {
					peny.BebanPenyusutanSebelumAtribusi = float64(peny.BulanAtribusi-1) * prevPeny.BebanPenyusutanSetelahAtribusi
					if prevPeny.NilaiBukuTahunBerkenaan < peny.BebanPenyusutanSebelumAtribusi {
						peny.BebanPenyusutanSebelumAtribusi = prevPeny.NilaiBukuTahunBerkenaan
					}
				}
			} else {
				peny.BebanPenyusutanSebelumAtribusi = 0
			}
		}
	}

	peny.BebanPenyusutanSebelumAtribusi = math.Ceil(peny.BebanPenyusutanSebelumAtribusi)

	if peny.BebanPenyusutanSebelumAtribusi < 0 {
		peny.BebanPenyusutanSebelumAtribusi = 0
	}

}
