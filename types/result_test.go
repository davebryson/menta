package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// More as an example of use than for testing functionalitys
func TestResult(t *testing.T) {
	assert := assert.New(t)

	// Ok response with tags
	// TODO: Upgrade to events
	/*r1 := Result{}
	r1.Tags = Tags{
		Tag{Key: []byte("dave"), Value: []byte("v")},
		Tag{Key: []byte("two"), Value: []byte("v")},
	}
	assert.Equal(2, len(r1.Tags))
	assert.Equal([]byte("dave"), r1.Tags[0].Key)*/

	// Error response
	bad := ErrorNoHandler()
	assert.Equal(uint32(1), bad.Code)
	assert.Equal("Handler not found", bad.Log)
}
