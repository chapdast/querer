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

	t.Run("SELECT_Variadic", func(t *testing.T) {
		Convey(t.Name(), t, func() {
			query, values, err := q.Build(
				SelectFields(fields[1],
					Coalesce(fields[0], "hello"),
				),
				Limit(100),
				Where(Conditional{
					Field:     fields[0],
					Operation: OprGreater,
					Value:     1,
				}),
				Where(Conditional{
					Field:     fields[0],
					Operation: OprNotEqual,
					Value:     120,
				}),
				Where(Conditional{
					Field:     fields[2],
					Operation: OprSubstring,
					Value:     "hello",
				}),
			)
			So(err, ShouldBeNil)
			So(query, ShouldEqual,
				"SELECT field_2, COALESCE(field_1, 'hello') FROM test WHERE field_1>$1 AND field_1!=$2 AND field_3 like '%'||$3||'%' LIMIT $4;")
			t.Log(values)
			So(values, ShouldHaveLength, 4)
			So(values[0], ShouldEqual, 1)
			So(values[1], ShouldEqual, 120)
			So(values[2], ShouldEqual, "hello")
			So(values[3], ShouldEqual, 100)
			t.Log(query)
			t.Log(values)
		})
	})

}
