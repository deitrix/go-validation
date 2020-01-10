package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstraints(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		Validate(nil, Constraints{
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Equal(t, 4, testConstraint.Calls)
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		testConstraint := &TestConstraint{}
		violations := Validate(nil, Constraints{
			testConstraint,
			testConstraint,
			testConstraint,
			testConstraint,
		})

		assert.Len(t, violations, 4)
	})
}

func TestElements(t *testing.T) {
	t.Run("should run all constraints", func(t *testing.T) {
		t.Run("against an array of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := [6]int{1, 2, 3, 4, 5, 6}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})

		t.Run("against a map of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10, 11: 12}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})

		t.Run("against a slice of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := []int{1, 2, 3, 4, 5, 6}

			Validate(values, Elements{
				testConstraint,
			})

			assert.Equal(t, 6, testConstraint.Calls)
		})
	})

	t.Run("should return all constraint violations", func(t *testing.T) {
		t.Run("against an array of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := [6]int{1, 2, 3, 4, 5, 6}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})

		t.Run("against a map of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10, 11: 12}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})

		t.Run("against a slice of values", func(t *testing.T) {
			testConstraint := &TestConstraint{}
			values := []int{1, 2, 3, 4, 5, 6}

			violations := Validate(values, Elements{
				testConstraint,
			})

			assert.Len(t, violations, 6)
		})
	})

	t.Run("should return no violations if the given value is nil", func(t *testing.T) {
		// NOTE: An array cannot be nil, and must have the length specified by it's type.

		t.Run("against a map", func(t *testing.T) {
			var m map[string]interface{}

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})

		t.Run("against a slice", func(t *testing.T) {
			var s []string

			testConstraint := &TestConstraint{}
			violations := Validate(s, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)

		})
	})

	t.Run("should return no violations if the given value's length is 0", func(t *testing.T) {
		// NOTE: An array cannot be nil, and must have the length specified by it's type.

		t.Run("against a map", func(t *testing.T) {
			m := make(map[string]interface{}, 0)

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)
		})

		t.Run("against a slice", func(t *testing.T) {
			s := make([]string, 0)

			testConstraint := &TestConstraint{}
			violations := Validate(s, Elements{
				testConstraint,
			})

			assert.Equal(t, 0, testConstraint.Calls)
			assert.Len(t, violations, 0)

		})
	})

	t.Run("should update the context's value node to the elements of the given value", func(t *testing.T) {
		t.Run("against an array", func(t *testing.T) {
			a := [2]string{"Hello", "World"}

			var values []string

			Validate(a, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					values = append(values, ctx.Value().Node.Interface().(string))
					return nil
				}),
			})

			require.Len(t, values, len(a))
			for i := range a {
				assert.Equal(t, a[i], values[i])
			}
		})

		t.Run("against a map", func(t *testing.T) {
			m := map[string]string{"Hello": "World"}

			var value string
			Validate(m, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					value = ctx.Value().Node.Interface().(string)
					return nil
				}),
			})

			assert.Equal(t, m["Hello"], value)
		})

		t.Run("against a slice", func(t *testing.T) {
			s := []string{"Hello", "World"}

			var values []string

			Validate(s, Elements{
				ConstraintFunc(func(ctx Context) []ConstraintViolation {
					values = append(values, ctx.Value().Node.Interface().(string))
					return nil
				}),
			})

			require.Len(t, values, len(s))
			for i := range s {
				assert.Equal(t, s[i], values[i])
			}
		})
	})

	t.Run("should update the path", func(t *testing.T) {
		t.Run("against an array", func(t *testing.T) {
			a := [2]string{"Hello", "World"}

			testConstraint := &TestConstraint{}
			violations := Validate(a, Elements{
				testConstraint,
			})

			require.Len(t, violations, 2)
			assert.Equal(t, ".[0]", violations[0].Path)
			assert.Equal(t, ".[1]", violations[1].Path)
		})

		t.Run("against a map", func(t *testing.T) {
			m := map[string]string{
				"Hello": "World",
			}

			testConstraint := &TestConstraint{}
			violations := Validate(m, Elements{
				testConstraint,
			})

			require.Len(t, violations, 1)
			assert.Equal(t, ".Hello", violations[0].Path)
		})

		t.Run("against a slice", func(t *testing.T) {
			a := []string{"Hello", "World"}

			testConstraint := &TestConstraint{}
			violations := Validate(a, Elements{
				testConstraint,
			})

			require.Len(t, violations, 2)
			assert.Equal(t, ".[0]", violations[0].Path)
			assert.Equal(t, ".[1]", violations[1].Path)
		})
	})
}

func TestFields(t *testing.T) {

}

func TestKeys(t *testing.T) {

}

func TestLazy(t *testing.T) {

}

func TestMap(t *testing.T) {

}

func TestWhen(t *testing.T) {

}

func TestWhenFn(t *testing.T) {

}

type TestConstraint struct {
	Calls int
}

func (c *TestConstraint) Violations(ctx Context) []ConstraintViolation {
	c.Calls++
	return []ConstraintViolation{
		ctx.Violation("test violations", nil),
	}
}
