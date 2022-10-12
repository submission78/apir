package main

// Test suite for the single-server VPIR scheme

import (
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/cloudflare/circl/group"
	"github.com/submission78/apir/lib/client"
	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/monitor"
	"github.com/submission78/apir/lib/server"
	"github.com/submission78/apir/lib/utils"
	"github.com/stretchr/testify/require"
)

func TestDH(t *testing.T) {
	dbLen := 1024 * 1024 // dbLen is specified in bits
	dbPRG := utils.RandomPRG()
	ecg := group.P256
	db := database.CreateRandomEllipticWithDigest(dbPRG, dbLen, ecg, true)
	fmt.Println("DB created")
	prg := utils.RandomPRG()
	retrieveBlocksDH(t, prg, db, "Diffie-Hellman")
}

func retrieveBlocksDH(t *testing.T, rnd io.Reader, db *database.Elliptic, testName string) {
	c := client.NewDH(rnd, &db.Info)
	s := server.NewDH(db)

	var i int
	totalTimer := monitor.NewMonitor()
	for j := 0; j < 10; j++ {
		i = rand.Intn(db.NumRows * db.NumColumns)
		query, err := c.QueryBytes(i)
		require.NoError(t, err)

		a, err := s.AnswerBytes(query)
		require.NoError(t, err)

		res, err := c.ReconstructBytes(a)
		require.NoError(t, err)
		require.Equal(t, db.Entries[i], res)
	}
	fmt.Printf("\nTotalCPU time %s: %.1fms\n", testName, totalTimer.Record())
}
