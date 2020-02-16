package constraints

import (
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/seeruk/go-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeAfter(t *testing.T) {
	past := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
	present := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	future := time.Date(3000, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("should return no violations if the context's time is after the constraints 'after' time", func(t *testing.T) {
		violations := TimeAfter(past)(validation.NewContext(present))
		assert.Len(t, violations, 0)
	})

	t.Run("should return a violation if the context's time is not after the constraints 'after' time", func(t *testing.T) {
		violations := TimeAfter(future)(validation.NewContext(present))
		assert.Len(t, violations, 1)
	})

	t.Run("should be optional (i.e. only applied if value is not empty)", func(t *testing.T) {
		violations := TimeAfter(past)(validation.NewContext(time.Time{}))
		assert.Len(t, violations, 0)
	})

	t.Run("should return details about the time the value should be after with a violation", func(t *testing.T) {
		violations := TimeAfter(future)(validation.NewContext(present))
		require.Len(t, violations, 1)
		assert.Equal(t, map[string]interface{}{
			"time": future.Format(time.RFC3339),
		}, violations[0].Details)
	})

	t.Run("should not panic if given a nil pointer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			TimeAfter(past)(validation.NewContext((*time.Time)(nil)))
		})
	})

	t.Run("should panic if given a value of the wrong type, even if it's empty", func(t *testing.T) {
		assert.Panics(t, func() { TimeAfter(past)(validation.NewContext("test")) })
		assert.Panics(t, func() { TimeAfter(past)(validation.NewContext(123)) })
		assert.Panics(t, func() { TimeAfter(past)(validation.NewContext(url.Values{})) })
		assert.Panics(t, func() { TimeAfter(past)(validation.NewContext(regexp.MustCompile("^test"))) })
	})

	t.Run("should return violations if given a value of the wrong type, even if it's empty, if strict types is false", func(t *testing.T) {
		ctx1 := validation.NewContext("")
		ctx1.StrictTypes = false
		ctx2 := validation.NewContext(0)
		ctx2.StrictTypes = false
		ctx3 := validation.NewContext(url.Values{})
		ctx3.StrictTypes = false

		assert.Len(t, TimeAfter(past)(ctx1), 1)
		assert.Len(t, TimeAfter(past)(ctx2), 1)
		assert.Len(t, TimeAfter(past)(ctx3), 1)
	})
}
