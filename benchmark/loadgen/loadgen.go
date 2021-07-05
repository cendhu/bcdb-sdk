package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/IBM-Blockchain/bcdb-sdk/benchmark/loadgen/writeonly"
	"github.com/IBM-Blockchain/bcdb-sdk/pkg/bcdb"
	"github.com/IBM-Blockchain/bcdb-sdk/pkg/config"
	"github.com/IBM-Blockchain/bcdb-server/pkg/logger"
)

func main() {
	c, err := ReadConfig("./config.yml")
	if err != nil {
		fmt.Errorf(err.Error())
	}

	logger, err := logger.New(
		&logger.Config{
			Level:         c.ConnectionConfig.LogLevel,
			OutputPath:    []string{"stdout"},
			ErrOutputPath: []string{"stderr"},
			Encoding:      "console",
			Name:          "bcdb-client",
		},
	)

	conConf := &config.ConnectionConfig{
		ReplicaSet: c.ConnectionConfig.ReplicaSet,
		RootCAs:    c.ConnectionConfig.RootCAs,
		Logger:     logger,
	}

	db, err := bcdb.Create(conConf)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	session, err := db.Session(&c.SessionConfig)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	totalKeys := c.WorkloadConfig.LoadPerClient * int(c.WorkloadConfig.Runtime.Seconds()) * c.WorkloadConfig.NumOfClients
	keysPerClient := totalKeys / c.WorkloadConfig.NumOfClients

	var clients []*writeonly.Client
	startKey := 1
	endKey := keysPerClient
	for i := 0; i < c.WorkloadConfig.NumOfClients; i++ {
		client := writeonly.New(i, session, "key", startKey, endKey, c.WorkloadConfig.LoadPerClient)
		clients = append(clients, client)
		startKey = endKey + 1
		endKey += keysPerClient
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(clients))
	start := time.Now()
	for _, c := range clients {
		go c.Run(wg)
	}
	wg.Wait()
	fmt.Printf("to write %d kv pairs (1 per tx), it took around %s msecs\n", totalKeys, time.Since(start).String())
}
