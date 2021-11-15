package integron

type Persister interface {
	Save([]Integron) error
}
