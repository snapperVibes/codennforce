package s_test

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/snapperVibes/codennforce"
)

func TestParse(t *testing.T) {
	ensure.DeepEqual(t, codennforce.ParseName([]byte("CONNELLY ELIZA     ")), "CONNELLY ELIZA")
}
