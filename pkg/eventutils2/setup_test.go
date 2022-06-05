package eventutils2

import (
	"os"
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
)

func TestMain(m *testing.M) {
	utils.Time = clock.NewMock()
	code := m.Run()
	os.Exit(code)
}
