package osexitanalyzer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitAnalyzer(t *testing.T) {
	path, err := filepath.Abs("../../../testdata")
	require.NoError(t, err, "error finding test data path")
	analysistest.Run(t, path, NewOsExitAnalyzer(), "./...")
}
