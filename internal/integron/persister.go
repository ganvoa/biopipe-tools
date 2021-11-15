package integron

type IntegronPersister interface {
	Save([]Integron) error
}
