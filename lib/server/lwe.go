package server

import (
	"github.com/submission78/apir/lib/database"
	"github.com/submission78/apir/lib/matrix"
)

type LWE struct {
	db *database.LWE
}

func NewLWE(db *database.LWE) *LWE {
	return &LWE{db: db}
}

func (s *LWE) DBInfo() *database.Info {
	return &s.db.Info
}

func (s *LWE) AnswerBytes(q []byte) ([]byte, error) {
	a := s.Answer(matrix.BytesToMatrix(q))
	return matrix.MatrixToBytes(a), nil
}

// Answer function for the LWE-based scheme. The query is represented as a
// vector
func (s *LWE) Answer(q *matrix.Matrix) *matrix.Matrix {
	return matrix.BinaryMul(q, s.db.Matrix)
}
