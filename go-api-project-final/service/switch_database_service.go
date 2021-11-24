package service

import "github.com/mandarinkb/go-api-project-final/repository"

type switchDatabaseService struct {
	swDatabaseRepo repository.SwitchDatabaseRepository
}

func NewSwDatabaseService(swDatabaseRepo repository.SwitchDatabaseRepository) SwitchDatabaseService {
	return switchDatabaseService{swDatabaseRepo}
}

func (sw switchDatabaseService) Read() (swDb []SwitchDatabase, err error) {
	swDatabaseRepo, err := sw.swDatabaseRepo.Read()
	if err != nil {
		return nil, err
	}
	for _, row := range swDatabaseRepo {
		swDb = append(swDb, mapDataSwDbResponse(row))
	}
	return swDb, nil
}

func (sw switchDatabaseService) ReadById(id int) (*SwitchDatabase, error) {
	swDbRepo, err := sw.swDatabaseRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	swDbRes := mapDataSwDbResponse(*swDbRepo)
	return &swDbRes, nil
}

func (sw switchDatabaseService) Create(swDb SwitchDatabase) error {
	return sw.swDatabaseRepo.Create(mapDataSwDbRequest(swDb))
}

func (sw switchDatabaseService) Update(swDb SwitchDatabase) error {
	return sw.swDatabaseRepo.Update(mapDataSwDbRequest(swDb))
}

func (sw switchDatabaseService) UpdateStatus(swDb SwitchDatabase) error {
	return sw.swDatabaseRepo.UpdateStatus(mapDataSwDbRequest(swDb))
}

func (sw switchDatabaseService) Delete(id int) error {
	return sw.swDatabaseRepo.Delete(id)
}

// แปลงค่า เพื่อส่งไปยัง repository
func mapDataSwDbRequest(swDb SwitchDatabase) repository.SwitchDatabase {
	return repository.SwitchDatabase{
		DatabaseId:     swDb.DatabaseId,
		DatabaseName:   swDb.DatabaseName,
		DatabaseStatus: swDb.DatabaseStatus,
	}
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataSwDbResponse(swDbRepo repository.SwitchDatabase) SwitchDatabase {
	return SwitchDatabase{
		DatabaseId:     swDbRepo.DatabaseId,
		DatabaseName:   swDbRepo.DatabaseName,
		DatabaseStatus: swDbRepo.DatabaseStatus,
	}
}
