package writeonly

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/IBM-Blockchain/bcdb-sdk/pkg/bcdb"
)

const (
	db = "bdb"
)

type Client struct {
	id                    int
	session               bcdb.DBSession
	keyPrefix             string
	keyStart              int
	keyEnd                int
	transactionsPerSecond int
	txInputs              []*kv
}

type kv struct {
	key   string
	value []byte
}

func New(clientID int, session bcdb.DBSession, keyPrefix string, keyStartNum, keyEndNum, tps int) *Client {
	var txInputs []*kv
	for i := keyStartNum; i <= keyEndNum; i++ {
		b := make([]byte, 1024)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		txInputs = append(
			txInputs,
			&kv{
				key:   fmt.Sprintf(keyPrefix+"%d", i),
				value: b,
			},
		)
	}

	return &Client{
		id:                    clientID,
		session:               session,
		keyPrefix:             keyPrefix,
		keyStart:              keyStartNum,
		keyEnd:                keyEndNum,
		transactionsPerSecond: tps,
		txInputs:              txInputs,
	}
}

func (c *Client) Run(wg *sync.WaitGroup) {
	waitMsec := time.Duration(1000/c.transactionsPerSecond) * time.Millisecond

	for i, input := range c.txInputs {
		go c.performDataTx(i, input)
		time.Sleep(waitMsec)
	}

	wg.Done()
}

func (c *Client) performDataTx(i int, kv *kv) {
	start := time.Now()
	tx, err := c.session.DataTx()
	if err != nil {
		panic(err)
	}

	if err := tx.Put("bdb", kv.key, kv.value, nil); err != nil {
		panic(err)
	}

	_, _, err = tx.Commit(true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client: %d, Tx: %d, Respone time: %s\n", c.id, i, time.Since(start).String())
}
