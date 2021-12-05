package integron

type IntegronRepository interface {
	AddIntegron(strainId int, integrons []Integron) error
	FindWithoutIntegronResult() []Integron
}
