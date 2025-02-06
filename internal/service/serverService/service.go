package serverservice

import (
	db "templates_new/internal/client"
	"templates_new/internal/repository"
	"templates_new/internal/service"
)

type serv struct {
	serverRepository repository.ServiceRepository
	txManager        db.TxManager
}

func NewService(
	serverRepository repository.ServiceRepository,
	txManager db.TxManager,
) service.ServerService {
	return &serv{
		serverRepository: serverRepository,
		txManager:        txManager,
	}
}
