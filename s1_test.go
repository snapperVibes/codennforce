package s_test

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestParse(t *testing.T) {
	ensure.DeepEqual(t, ParseName([]byte("CONNELLY ELIZA     ")), "CONNELLY ELIZA")
}
