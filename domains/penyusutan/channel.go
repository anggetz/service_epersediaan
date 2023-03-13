package penyusutan

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Channel struct {
	c *nats.Conn
}

func NewChannel(c *nats.Conn) *Channel {
	return &Channel{
		c: c,
	}
}

type reqRegisterCalcPenyusutan struct {
	OpdId     int    `json:"opd_id"`
	JenisAset string `json:"jenis_aset"`
}

func (c *Channel) RegisterCalcPenyusutan() {
	c.c.Subscribe("simada.calcpenyusutan", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data), time.Now())

		var dataPayload reqRegisterCalcPenyusutan

		err := json.Unmarshal(msg.Data, &dataPayload)

		if err != nil {
			log.Println("Error", err.Error())
			c.c.Publish(msg.Reply, []byte("Error!"))
			return
		}

		NewUseCase().CalcPenyusutan(dataPayload.OpdId, dataPayload.JenisAset, time.Date(2022, 12, 0, 0, 0, 0, 0, time.Local))

		log.Println("Done!")
		c.c.Publish(msg.Reply, []byte("Done!"))

	})
}
