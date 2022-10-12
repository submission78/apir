package server

import (
	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/matrix"
)

type Amplify struct {
	lwe *LWE
}

func NewAmplify(db *database.LWE) *Amplify {
	return &Amplify{
		lwe: NewLWE(db),
	}
}

func (a *Amplify) DBInfo() *database.Info {
	return &a.lwe.db.Info
}

func (a *Amplify) Answer(qq []*matrix.Matrix) []*matrix.Matrix {
	ans := make([]*matrix.Matrix, len(qq))
	for i, q := range qq {
		ans[i] = matrix.BinaryMul(q, a.lwe.db.Matrix)
	}

	return ans
}

func (a *Amplify) AnswerBytes(qq []byte) ([]byte, error) {
	ans := a.Answer(matrix.BytesToMatrices(qq))

	// encode
	return matrix.MatricesToBytes(ans), nil
}
