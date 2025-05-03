package logger_test

import (
	"testing"

	"github.com/Te8va/Gofermarch/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	str := "abc_tmp.txt"
	logger.SetLogFile(str)

	n, err := logger.Logger().Write([]byte(str))
	require.NoError(t, err)
	require.Equal(t, len(str), n)
}
