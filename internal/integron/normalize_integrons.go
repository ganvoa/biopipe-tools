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

				has_qac := false
				is_calin := false
				is_zero_one := false
				is_complete := false
				is_type := ""

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
						has_qac = true
					}

					normalizedSlice = strings.Split(normalized, " ")
					if normalizedSlice[0] == "intI" {
						if strings.Contains(normalizedSlice[len(normalizedSlice)-1], "qac") {
							is_complete = true
						}
					}

					if normalizedSlice[0] != "intI" {
						is_calin = true
					}

					inverted, err = in.Invert(integronSlice)
					if err != nil {
						in.logger.Warnf("error %v", err)
						continue
					}
				} else {
					if integronSlice[0] == "intI" {
						is_zero_one = true
					}
				}

				if is_zero_one {
					is_type = "zero_one"
				} else if is_complete {
					is_type = "complete"
				} else if is_calin {
					is_type = "calin"
				}

				in.logger.Debugf("original %s", integron)
				in.logger.Debugf("normalized %s", normalized)
				in.logger.Debugf("inverted %s", inverted)
				in.logger.Debugf("short %s", short)
				in.logger.Debugf("is_zero_one %b", is_zero_one)
				in.logger.Debugf("is_complete %b", is_complete)
				in.logger.Debugf("is_calin %b", is_calin)
				in.logger.Debugf("has_qac %b", has_qac)
				in.logger.Debugf("is_type %b", is_type)

				integronUpdate = append(integronUpdate, fasta.IntegronResult{
					Original:    integron,
					Normalized:  normalized,
					Inverted:    inverted,
					Short:       short,
					Is_Zero_One: is_zero_one,
					Is_Complete: is_complete,
					Is_Calin:    is_calin,
					Has_Qac:     has_qac,
					Is_Type:     is_type,
				})

				in.logger.Debug("-----------------------")

			}

			in.fastaRepository.UpdateIntegron(strain.Id, integronUpdate)
		}

	}

	return nil
}
