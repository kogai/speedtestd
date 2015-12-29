package speedtestd

import (
	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
}

func (service *Service) StartJob() (string, error) {
}

func (service *Service) StopJob() (string, error) {
}

func (service *Service) showResult() (string, error) {
}
