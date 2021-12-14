package integron

import (
	"strings"

	"github.com/ganvoa/biopipe-tools/internal"
	"github.com/ganvoa/biopipe-tools/internal/fasta"
)

type IntegronNormalizer struct {
	cleaner         IntegronResultCleaner
	parser          IntegronParser
	fastaRepository fasta.FastaRepository
	logger          internal.Logger
}

func NewIntegronNormalizer(cleaner IntegronResultCleaner, parser IntegronParser, repository fasta.FastaRepository, logger internal.Logger) IntegronNormalizer {
	in := IntegronNormalizer{}
	in.cleaner = cleaner
	in.parser = parser
	in.fastaRepository = repository
	in.logger = logger
	return in
}

func (in IntegronNormalizer) Invert(integron []string) (string, error) {
	delimiter := " "
	chain := ""

	for i := len(integron) - 1; i >= 0; i-- {
		if i > 0 {
			chain = chain + integron[i] + delimiter
		} else {
			chain = chain + integron[i]
		}

	}

	return chain, nil
}

func (in IntegronNormalizer) Run() error {

	currentFrom := 10000000000

	for {

		in.logger.Debug("finding strains with integrons")
		strains, err := in.fastaRepository.FindWithIntegronResult(currentFrom)

		if err != nil {
			return err
		}

		in.logger.Debugf("found %d strains", len(strains))
		if len(strains) == 0 {
			break
		}

		for _, strain := range strains {

			integronUpdate := []fasta.IntegronResult{}

			in.logger.Debug("-----------------------")
			in.logger.Debugf("strain %d", strain.Id)

			for _, integron := range strain.Integrons {

				in.logger.Debug("normalizando integron")

				integronSlice := strings.Split(integron, " ")
				normalized := integron
				inverted := integron
				short := integron

				in.logger.Debugf("first %s", integronSlice[0])
				in.logger.Debugf("last %s", integronSlice[len(integronSlice)-1])

				if len(integronSlice) > 1 {

					if integronSlice[len(integronSlice)-1] == "intI" || strings.Contains(integronSlice[0], "qac") {
						normalized, err = in.Invert(integronSlice)
						if err != nil {
							in.logger.Warnf("error %v", err)
							continue
						}
					}

					normalizedSlice := strings.Split(normalized, " ")
					if normalizedSlice[0] == "intI" {
						normalizedSlice = normalizedSlice[1:]
						short = strings.Join(normalizedSlice, " ")
					}

					if strings.Contains(normalizedSlice[len(normalizedSlice)-1], "qac") {
						short = strings.Join(normalizedSlice[:len(normalizedSlice)-1], " ")
					}

					inverted, err = in.Invert(integronSlice)
					if err != nil {
						in.logger.Warnf("error %v", err)
						continue
					}
				}

				in.logger.Debugf("original %s", integron)
				in.logger.Debugf("normalized %s", normalized)
				in.logger.Debugf("inverted %s", inverted)
				in.logger.Debugf("short %s", short)

				integronUpdate = append(integronUpdate, fasta.IntegronResult{
					Original:   integron,
					Normalized: normalized,
					Inverted:   inverted,
					Short:      short,
				})

				in.logger.Debug("-----------------------")

			}

			in.fastaRepository.UpdateIntegron(strain.Id, integronUpdate)
		}

	}

	return nil
}
