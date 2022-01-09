package querer

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNew(t *testing.T) {
	q := New("test", []string{"field_1", "field_2", "field_3"})
	t.Run("SELECT", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, err := q.Build(
				Select(),
				Where(map[string]OperatorType{
					"field_1": OprGreater,
				}),
			)
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"SELECT field_1, field_2, field_3 FROM test WHERE field_1>$1;")
			t.Log(query)
		})
	})
	t.Run("INSERT", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, err := q.Build(
				Insert(),
				Where(map[string]OperatorType{
					"field_1": OprGreater,
				}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"INSERT INTO test (field_1, field_2, field_3) VALUES ($1, $2, $3 );")
			t.Log(query)
		})
	})
	t.Run("UPDATE", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, err := q.Build(
				Update(),
				Where(map[string]OperatorType{
					"field_1": OprInArray,
				}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"UPDATE test SET field_1=$1, field_2=$2, field_3=$3 WHERE field_1 = ANY($4);")
			t.Log(query)
		})
	})
	t.Run("DELETE", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, err := q.Build(
				Delete(),
				Where(
					map[string]OperatorType{
						"field_1": OprInArray,
					}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"DELETE FROM test WHERE field_1 = ANY($1);")
			t.Log(query)
		})
	})

	t.Run("SELECT_II", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, err := q.Build(
				Select(),
				Where(
					map[string]OperatorType{
						"field_1": OprGreater,
						"field_2": OprArrayOverlap,
					}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"SELECT field_1, field_2, field_3 FROM test WHERE field_1>$1 AND field_2 && $2;")
			t.Log(query)
		})
	})
}
