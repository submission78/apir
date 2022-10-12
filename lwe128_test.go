package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/submission78/apir/lib/client"
	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/server"
	"github.com/submission78/apir/lib/utils"
	"github.com/stretchr/testify/require"
)

func TestLWE128(t *testing.T) {
	dbLen := 1024 * 1024 // dbLen is specified in bits
	db := database.CreateRandomBinaryLWEWithLength128(utils.RandomPRG(), dbLen)
	p := utils.ParamsWithDatabaseSize128(db.Info.NumRows, db.Info.NumColumns)
	retrieveBlocksLWE128(t, db, p, "TestLWE128")
}

func retrieveBlocksLWE128(t *testing.T, db *database.LWE128, params *utils.ParamsLWE, testName string) {
	c := client.NewLWE128(utils.RandomPRG(), &db.Info, params)
	s := server.NewLWE128(db)

	ti := time.Now()
	for j := 0; j < 100; j++ {
		i := rand.Intn(params.L * params.M)
		query, err := c.QueryBytes(i)
		require.NoError(t, err)

		a, err := s.AnswerBytes(query)
		require.NoError(t, err)

		res, err := c.ReconstructBytes(a)
		require.NoError(t, err)
		require.Equal(t, uint32(db.Matrix.Get(utils.VectorToMatrixIndices(i, db.Info.NumColumns))), res)
	}
	fmt.Printf("Total time %s: %.1fs\n", testName, time.Since(ti).Seconds())
}
