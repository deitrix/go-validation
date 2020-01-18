package proto

import (
	"errors"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapToStruct(t *testing.T) {
	t.Run("should return nil if an empty map is given", func(t *testing.T) {
		assert.Empty(t, MapToStruct(map[string]interface{}{}))
	})

	t.Run("should be able to handle all expected types", func(t *testing.T) {
		input := mapOfSupportedTypes()
		input["nested"] = mapOfSupportedTypes()

		output := MapToStruct(input)

		assertMapOfSupportedTypes(t, input, output)

		require.IsType(t, map[string]interface{}{}, input["nested"])
		require.IsType(t, &types.Value_StructValue{}, output.Fields["nested"].Kind)

		// Well, this is gross.
		assertMapOfSupportedTypes(t,
			input["nested"].(map[string]interface{}),
			output.Fields["nested"].Kind.(*types.Value_StructValue).StructValue,
		)
	})
}

func assertMapOfSupportedTypes(t *testing.T, input map[string]interface{}, output *types.Struct) {
	// nil
	assert.IsType(t, &types.Value_NullValue{}, output.Fields["nil"].Kind)
	assert.Equal(t, output.Fields["nil"], &types.Value{Kind: &types.Value_NullValue{
		NullValue: types.NullValue_NULL_VALUE,
	}})

	// bool
	require.IsType(t, true, input["bool"])
	require.IsType(t, &types.Value_BoolValue{}, output.Fields["bool"].Kind)
	assert.Equal(t, output.Fields["bool"], &types.Value{Kind: &types.Value_BoolValue{
		BoolValue: input["bool"].(bool),
	}})

	// int
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["int"].Kind)
	assert.Equal(t, output.Fields["int"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["int"].(int)),
	}})

	// int8
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["int8"].Kind)
	assert.Equal(t, output.Fields["int8"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["int8"].(int8)),
	}})

	// int16
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["int16"].Kind)
	assert.Equal(t, output.Fields["int16"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["int16"].(int16)),
	}})

	// int32
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["int32"].Kind)
	assert.Equal(t, output.Fields["int32"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["int32"].(int32)),
	}})

	// int64
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["int64"].Kind)
	assert.Equal(t, output.Fields["int64"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["int64"].(int64)),
	}})

	// uint
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["uint"].Kind)
	assert.Equal(t, output.Fields["uint"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["uint"].(uint)),
	}})

	// uint8
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["uint8"].Kind)
	assert.Equal(t, output.Fields["uint8"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["uint8"].(uint8)),
	}})

	// uint16
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["uint16"].Kind)
	assert.Equal(t, output.Fields["uint16"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["uint16"].(uint16)),
	}})

	// uint32
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["uint32"].Kind)
	assert.Equal(t, output.Fields["uint32"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["uint32"].(uint32)),
	}})

	// uint64
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["uint64"].Kind)
	assert.Equal(t, output.Fields["uint64"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["uint64"].(uint64)),
	}})

	// float32
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["float32"].Kind)
	assert.Equal(t, output.Fields["float32"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: float64(input["float32"].(float32)),
	}})

	// float64
	require.IsType(t, &types.Value_NumberValue{}, output.Fields["float64"].Kind)
	assert.Equal(t, output.Fields["float64"], &types.Value{Kind: &types.Value_NumberValue{
		NumberValue: input["float64"].(float64),
	}})

	// string
	require.IsType(t, &types.Value_StringValue{}, output.Fields["string"].Kind)
	assert.Equal(t, output.Fields["string"], &types.Value{Kind: &types.Value_StringValue{
		StringValue: input["string"].(string),
	}})

	// error
	require.IsType(t, &types.Value_StringValue{}, output.Fields["error"].Kind)
	assert.Equal(t, output.Fields["error"], &types.Value{Kind: &types.Value_StringValue{
		StringValue: input["error"].(error).Error(),
	}})

	// unsupported value
	require.IsType(t, &types.Value_StringValue{}, output.Fields["unsupported"].Kind)
	assert.Equal(t, output.Fields["unsupported"], &types.Value{Kind: &types.Value_StringValue{
		StringValue: input["unsupported"].(time.Duration).String(),
	}})
}

func mapOfSupportedTypes() map[string]interface{} {
	return map[string]interface{}{
		"nil":         nil,
		"bool":        true,
		"int":         123,
		"int8":        int8(123),
		"int16":       int16(123),
		"int32":       int32(123),
		"int64":       int64(123),
		"uint":        uint(123),
		"uint8":       uint8(123),
		"uint16":      uint16(123),
		"uint32":      uint32(123),
		"uint64":      uint64(123),
		"float32":     float32(123.456),
		"float64":     123.456,
		"string":      "test",
		"error":       errors.New("test error"),
		"unsupported": 15 * time.Second,
	}
}
