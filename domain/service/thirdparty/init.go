package thirdparty

import (
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/database"
)

type Service struct {
	Config config.Config            `inject:"config"`
	DB     *database.EngineDatabase `inject:"database"`
}

func (svc *Service) StartUp() {}

func (svc *Service) Shutdown() {}
