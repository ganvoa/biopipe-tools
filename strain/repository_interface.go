package strain

type StrainRepository interface {
	SaveAll(strains []Strain) error
}
