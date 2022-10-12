package client

import (
	"errors"
	"io"

	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/matrix"
	"github.com/submission78/apir/lib/utils"
)

// LEW based authenticated single server PIR client

// Client description
type LWE struct {
	dbInfo *database.Info
	state  *StateLWE
	params *utils.ParamsLWE
	rnd    io.Reader
}

type StateLWE struct {
	A      *matrix.Matrix
	digest *matrix.Matrix
	secret *matrix.Matrix
	i      int
	j      int
	t      uint32
}

func NewLWE(rnd io.Reader, info *database.Info, params *utils.ParamsLWE) *LWE {
	return &LWE{
		dbInfo: info,
		params: params,
		rnd:    rnd,
	}
}

func (c *LWE) Query(i, j int) *matrix.Matrix {
	// Lazy way to sample a random scalar
	rand := matrix.NewRandom(c.rnd, 1, 1)

	// digest is already stored in the state when receiving the database info
	c.state = &StateLWE{
		A:      matrix.NewRandom(utils.NewPRG(c.params.SeedA), c.params.N, c.params.L),
		digest: c.dbInfo.DigestLWE,
		secret: matrix.NewRandom(c.rnd, 1, c.params.N),
		i:      i,
		j:      j,
		t:      rand.Get(0, 0),
	}

	// Query has dimension 1 x l
	query := matrix.Mul(c.state.secret, c.state.A)

	// Error has dimension 1 x l
	e := matrix.NewGauss(1, c.params.L)

	msg := matrix.New(1, c.params.L)
	msg.Set(0, i, c.state.t)

	query.Add(e)
	query.Add(msg)

	return query
}

func (c *LWE) QueryBytes(index int) ([]byte, error) {
	i, j := utils.VectorToMatrixIndices(index, c.dbInfo.NumColumns)
	m := c.Query(i, j)
	return matrix.MatrixToBytes(m), nil
}

func (c *LWE) Reconstruct(answers *matrix.Matrix) (uint32, error) {
	s_trans_d := matrix.Mul(c.state.secret, c.state.digest)
	answers.Sub(s_trans_d)

	outs := make([]uint32, c.params.M)
	for i := 0; i < c.params.M; i++ {
		v := answers.Get(0, i)
		if c.inRange(v) {
			outs[i] = 0
		} else if c.inRange(v - c.state.t) {
			outs[i] = 1
		} else {
			return 0, errors.New("REJECT")
		}
	}

	return outs[c.state.j], nil
}

func (c *LWE) ReconstructBytes(a []byte) (uint32, error) {
	return c.Reconstruct(matrix.BytesToMatrix(a))
}

func (c *LWE) inRange(val uint32) bool {
	return (val < c.params.B) || (val > -c.params.B)
}
