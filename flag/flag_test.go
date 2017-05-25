package flag

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func Test_NewRequiredFlagError_ReturnsError_InCaseOfSuccess(t *testing.T) {
	// Act
	err := NewRequiredFlagError("something")

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify something")
}

func Test_Get_ReturnsFlagValue_InCaseOfFlagSet(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("foo", "", "")
	if err := set.Parse([]string{"--foo", "bar"}); err != nil {
		assert.FailNow(t, err.Error())
	}
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
