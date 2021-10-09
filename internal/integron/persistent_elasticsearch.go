package integron

import (
	"github.com/ganvoa/biopipe-tools/internal"
)

type integronPersistentES struct {
	strainId   int
	repository IntegronRepository
	logger     internal.Logger
}

func NewIntegronPersistentES(strainId int, repository IntegronRepository, logger internal.Logger) integronPersistentES {
	ipf := integronPersistentES{
		repository: repository,
		logger:     logger,
		strainId:   strainId,
	}
	return ipf
}

func (ipf integronPersistentES) Save(integrons []Integron) error {

	ipf.logger.Info("saving to elasticsearch")
	err := ipf.repository.AddIntegron(ipf.strainId, integrons)
	if err != nil {
		return err
	}

	return nil
}
