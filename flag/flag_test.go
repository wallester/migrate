package flag

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func Test_NewRequiredFlagError_ReturnsError_InCaseOfSuccess(t *testing.T) {
	// Act
	err := NewRequiredFlagError("something")

	// Assert
	assert.EqualError(t, err, "please specify something")
}

func Test_NewWrongFormatFlagError_ReturnsError_InCaseOfSuccess(t *testing.T) {
	// Act
	err := NewWrongFormatFlagError("something")

	// Assert
	assert.EqualError(t, err, "parsing something failed")
}

func Test_Get_ReturnsFlagValue_InCaseOfFlagSet(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("foo", "", "")
	require.NoError(t, set.Parse([]string{"--foo", "bar"}))

	c := cli.NewContext(nil, set, nil)

	// Act
	value := Get(c, "foo")

	// Assert
	assert.Equal(t, "bar", value)
}

func Test_Get_ReturnsEmptyString_InCaseOfFlagNotSet(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// Act
	value := Get(c, "foo")

	// Assert
	assert.Empty(t, value)
}

func Test_GetBool_ReturnsTrue_InCaseOfFlagSetWithValue(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.Bool("foo", true, "")
	require.NoError(t, set.Parse([]string{"--foo"}))

	c := cli.NewContext(nil, set, nil)

	// Act
	res := GetBool(c, "foo")

	// Assert
	assert.True(t, res)
}

func Test_GetBool_ReturnsFalse_InCaseOfFlagNotSet(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// Act
	res := GetBool(c, "foo")

	// Assert
	assert.False(t, res)
}
