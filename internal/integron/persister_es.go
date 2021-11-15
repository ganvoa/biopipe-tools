package integron

import (
	"github.com/ganvoa/biopipe-tools/internal"
)

type ElasticSearchPersister struct {
	strainId   int
	repository IntegronRepository
	logger     internal.Logger
}

func NewElasticSearchPersister(strainId int, repository IntegronRepository, logger internal.Logger) ElasticSearchPersister {
	ipf := ElasticSearchPersister{
		repository: repository,
		logger:     logger,
		strainId:   strainId,
	}
	return ipf
}

func (ipf ElasticSearchPersister) Save(integrons []Integron) error {

	ipf.logger.Info("saving to elasticsearch")
	err := ipf.repository.AddIntegron(ipf.strainId, integrons)
	if err != nil {
		return err
	}

	return nil
}
