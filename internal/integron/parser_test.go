package integron_test

import (
	"path/filepath"
	"testing"

	"github.com/ganvoa/biopipe-tools/internal/integron"
	"github.com/ganvoa/biopipe-tools/internal/platform"
	"github.com/stretchr/testify/assert"
)

var integronSuite = []struct {
	FolderName                string
	NumberOfIntegronsExpected int
	ResultsExpected           []string
}{
	{
		"Results_Integron_Finder_488901",
		2,
		[]string{
			"intI trim_DfrA1_like-NCBIFAM attC ANT_3pp_I-NCBIFAM attC SMR_qac_E-NCBIFAM",
			"attC attC protein protein",
		},
	},
	{
		"Results_Integron_Finder_488112",
		3,
		[]string{
			"intI trim_DfrA12-NCBIFAM attC attC ANT_3pp_I-NCBIFAM attC SMR_qac_E-NCBIFAM",
			"attC blaOXA-1_like-NCBIFAM attC AAC_6p_Ib-NCBIFAM",
			"attC attC protein protein",
		},
	},
}

func TestWhenFolderHasNoIntegronFilesShouldReturnZeroLengthList(t *testing.T) {
	folder := filepath.Join(".", "parser_test_data", "Results_Integron_Finder_488113")

	logger := &platform.FakeLogger{}
	parser := integron.NewParser(logger)
	integrons, err := parser.Parse(folder)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(integrons))
}

func TestWhenFolderHasIntegronFilesShouldReturnValidString(t *testing.T) {
	logger := &platform.FakeLogger{}

	for _, tt := range integronSuite {
		t.Run(tt.FolderName, func(t *testing.T) {
			folder := filepath.Join(".", "parser_test_data", tt.FolderName)
			parser := integron.NewParser(logger)
			integrons, err := parser.Parse(folder)
			assert.NoError(t, err)
			assert.Equal(t, tt.NumberOfIntegronsExpected, len(integrons))
			assert.Equal(t, tt.ResultsExpected, integrons)
		})
	}

}
