package querer

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNew(t *testing.T) {
	q := New("test")
	fields := []string{"field_1", "field_2", "field_3"}
	t.Run("SELECT", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, values, err := q.Build(
				Select(fields),
				Where(Conditional{
					Field:     fields[0],
					Operation: OprGreater,
					Value:     1,
				}),
			)
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"SELECT field_1, field_2, field_3 FROM test WHERE field_1>$1;")
			So(values[0], ShouldEqual, 1)
			t.Log(query)
			t.Log(values)
		})
	})
	t.Run("INSERT", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, values, err := q.Build(
				Insert(fields, []interface{}{
					10,
					"hello",
					[]uint64{4, 5, 6},
				}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"INSERT INTO test (field_1, field_2, field_3) VALUES ($1, $2, $3 );")
			So(values, ShouldHaveLength, 3)
			So(values[0], ShouldEqual, 10)
			So(values[1], ShouldEqual, "hello")
			So(values[2], ShouldResemble, []uint64{4, 5, 6})

			t.Log(query)
			t.Log(values)
		})
	})
	t.Run("UPDATE", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, values, err := q.Build(
				Update(fields, []interface{}{
					10,
					"hello",
					[]uint64{4, 5, 6},
				}),
				Where(Conditional{
					Field:     fields[0],
					Operation: OprInArray,
					Value:     []uint64{5, 11, 62, 64},
				}),
				Where(Conditional{
					Field:     fields[2],
					Operation: OprEqual,
					Value:     13,
				}),
			)
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"UPDATE test SET field_1=$1, field_2=$2, field_3=$3 WHERE field_1 = ANY($4) AND field_3=$5;")
			So(values, ShouldHaveLength, 5)
			So(values[0], ShouldEqual, 10)
			So(values[1], ShouldEqual, "hello")
			So(values[2], ShouldResemble, []uint64{4, 5, 6})
			So(values[3], ShouldResemble, []uint64{5, 11, 62, 64})
			So(values[4], ShouldEqual, 13)

			t.Log(query)
			t.Log(values)
		})
	})
	t.Run("DELETE", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, values, err := q.Build(
				Delete(),
				Where(Conditional{
					Field:     "field_3",
					Operation: OprInArray,
					Value:     []int{3, 5, 6, 7},
				}))
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"DELETE FROM test WHERE field_3 = ANY($1);")
			So(values, ShouldHaveLength, 1)
			So(values[0], ShouldResemble, []int{3, 5, 6, 7})

			t.Log(query)
			t.Log(values)
		})
	})
}
