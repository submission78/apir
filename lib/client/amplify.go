package client

import (
	"errors"
	"io"

	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/ecc"
	"github.com/submission78/apir/lib/matrix"
	"github.com/submission78/apir/lib/utils"
)

type Amplify struct {
	repetitions int    // 2*t + 1
	lwes        []*LWE // base client to each element of output of ECC
}

func NewAmplify(rnd io.Reader, info *database.Info, params *utils.ParamsLWE, tECC int) *Amplify {
	repetitions := tECC*2 + 1

	lwes := make([]*LWE, repetitions)
	for i := range lwes {
		lwes[i] = NewLWE(rnd, info, params)
	}

	return &Amplify{
		repetitions: repetitions,
		lwes:        lwes,
	}
}

func (a *Amplify) Query(i, j int) []*matrix.Matrix {
	queries := make([]*matrix.Matrix, a.repetitions)
	for k := 0; k < a.repetitions; k++ {
		queries[k] = a.lwes[k].Query(i, j)
	}
	return queries
}

func (a *Amplify) QueryBytes(index int) ([]byte, error) {
	i, j := utils.VectorToMatrixIndices(index, a.lwes[0].dbInfo.NumColumns)
	ms := a.Query(i, j)

	// encode
	return matrix.MatricesToBytes(ms), nil
}

func (a *Amplify) Reconstruct(answers []*matrix.Matrix) (uint32, error) {
	outputs := make([]uint32, a.repetitions)
	var err error
	for i := range outputs {
		outputs[i], err = a.lwes[i].Reconstruct(answers[i])
		if err != nil {
			return 0, errors.New("REJECT")
		}
	}

	// find and return majority
	ecc := ecc.New((a.repetitions - 1) / 2)
	return ecc.Decode(outputs)
}

func (a *Amplify) ReconstructBytes(answers []byte) (uint32, error) {
	return a.Reconstruct(matrix.BytesToMatrices(answers))
}
