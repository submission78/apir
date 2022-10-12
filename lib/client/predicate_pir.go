package client

import (
	"io"

	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/fss"
	"github.com/submission78/apir/lib/query"
)

// PredicatePIR represent the client for the FSS-based complex-queries non-verifiable PIR
type PredicatePIR struct {
	*clientFSS
}

// NewPredicatePIR returns a new client for the DPF-base multi-bit classical PIR
// scheme
func NewPredicatePIR(rnd io.Reader, info *database.Info) *PredicatePIR {
	executions := 1
	return &PredicatePIR{
		&clientFSS{
			rnd:        rnd,
			dbInfo:     info,
			state:      nil,
			Fss:        fss.ClientInitialize(executions), // only one value
			executions: executions,
		},
	}
}

// QueryBytes executes Query and encodes the result a byte array for each
// server
func (c *PredicatePIR) QueryBytes(in []byte, numServers int) ([][]byte, error) {
	return c.clientFSS.queryBytes(in, numServers)
}

// Query outputs the queries, i.e. DPF keys, for index i. The DPF
// implementation assumes two servers.
func (c *PredicatePIR) Query(q *query.ClientFSS, numServers int) []*query.FSS {
	return c.query(q, numServers)
}

// ReconstructBytes returns []byte
func (c *PredicatePIR) ReconstructBytes(answers [][]byte) (interface{}, error) {
	return c.reconstructBytes(answers)
}

// Reconstruct reconstruct the entry of the database from answers
func (c *PredicatePIR) Reconstruct(answers [][]uint32) (uint32, error) {
	return c.reconstruct(answers)
}
