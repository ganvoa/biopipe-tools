package integron

import (
	"fmt"
	"os/exec"

	"github.com/ganvoa/biopipe-tools/internal"
)

type integronFinder struct {
	outputDir string
	logger    internal.Logger
}

func NewIntegronFinder(outputDir string, logger internal.Logger) integronFinder {
	ifind := integronFinder{}
	ifind.outputDir = outputDir
	ifind.logger = logger
	return ifind
}

func (ifind integronFinder) Run(fastaFilePath string) error {
	// docker run --rm -it -v $(pwd)/download:/fasta -v $(pwd)/test/integrons:/output gamboa/biopipe:latest /bin/bash -c "source activate abricate; integron_finder --local-max --outdir /output --split-results --func-annot /fasta/488112.fasta"
	command := fmt.Sprintf(`docker run --rm -it -v "$(pwd)/download:/fasta" -v "$(pwd)/test/integrons:/output" gamboa/biopipe:latest /bin/bash -c "source activate abricate; integron_finder --local-max --outdir %s --split-results --func-annot %s`, ifind.outputDir, fastaFilePath)
	out, err := exec.Command(command).Output()

	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	return nil
}
