package pmx_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/wcamarao/pmx"
	"github.com/wcamarao/pmx/test"
)

type SelectSliceSuite struct {
	suite.Suite
	conn *pgx.Conn
}

func (s *SelectSliceSuite) SetupTest() {
	s.conn = test.Connect()
}

func TestSelectSlice(t *testing.T) {
	suite.Run(t, new(SelectSliceSuite))
}

func (s *SelectSliceSuite) TestPointer() {
	var samples []*test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select $1 as id, $2 as label", "a", "b")
	s.Equal([]*test.Sample{{ID: "a", Label: "b"}}, samples)
	s.Nil(err)
	s.True(ok)
}

func (s *SelectSliceSuite) TestSkipNull() {
	var samples []*test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select $1 as id, null as label", "a")
	s.Equal([]*test.Sample{{ID: "a"}}, samples)
	s.Nil(err)
	s.True(ok)
}

func (s *SelectSliceSuite) TestSkipTransient() {
	var samples []*test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select 'a' as id, 'b' as transient")
	s.Equal([]*test.Sample{{ID: "a"}}, samples)
	s.Nil(err)
	s.True(ok)
}

func (s *SelectSliceSuite) TestNoRows() {
	var samples []*test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select 1 limit 0")
	s.Empty(samples)
	s.Nil(err)
	s.False(ok)
}

func (s *SelectSliceSuite) TestValue() {
	var samples []*test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, samples, "select 1")
	s.Equal(pmx.ErrInvalidRef, err)
	s.False(ok)
}

func (s *SelectSliceSuite) TestPointerOfStructValue() {
	var samples []test.Sample
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select 1")
	s.Equal(pmx.ErrInvalidRef, err)
	s.False(ok)
}

func (s *SelectSliceSuite) TestPointerOfMapPointer() {
	var samples []*map[string]string
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select 1")
	s.Equal(pmx.ErrInvalidRef, err)
	s.False(ok)
}

func (s *SelectSliceSuite) TestPointerOfMapValue() {
	var samples []map[string]string
	ok, err := pmx.Select(context.Background(), s.conn, &samples, "select 1")
	s.Equal(pmx.ErrInvalidRef, err)
	s.False(ok)
}
