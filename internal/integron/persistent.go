package integron

type IntegronPersistent interface {
	Save([]Integron) error
}
