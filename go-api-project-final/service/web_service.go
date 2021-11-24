package service

import "github.com/mandarinkb/go-api-project-final/repository"

type webService struct {
	webServ repository.WebRepository
}

func NewWebService(webServ repository.WebRepository) WebService {
	return webService{webServ}
}

func (w webService) Read() (web []Web, err error) {
	webRepo, err := w.webServ.Read()
	if err != nil {
		return nil, err
	}
	for _, row := range webRepo {
		web = append(web, mapDataWebResponse(row))
	}

	return web, nil
}

func (w webService) ReadById(id int) (*Web, error) {
	webRepo, err := w.webServ.ReadById(id)
	if err != nil {
		return nil, err
	}
	webRes := mapDataWebResponse(*webRepo)
	return &webRes, nil
}

func (w webService) Create(web Web) error {
	return w.webServ.Create(mapDataWebRequest(web))
}

func (w webService) Update(web Web) error {
	return w.webServ.Update(mapDataWebRequest(web))
}

func (w webService) UpdateStatus(web Web) error {
	return w.webServ.UpdateStatus(mapDataWebRequest(web))
}

func (w webService) Delete(id int) error {
	return w.webServ.Delete(id)
}

// แปลงค่า เพื่อส่งไปยัง repository
func mapDataWebRequest(web Web) repository.Web {
	return repository.Web{
		WebId:     web.WebId,
		WebName:   web.WebName,
		WebUrl:    web.WebUrl,
		WebStatus: web.WebStatus,
		IconUrl:   web.IconUrl,
	}
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataWebResponse(webRepo repository.Web) Web {
	return Web{
		WebId:     webRepo.WebId,
		WebName:   webRepo.WebName,
		WebUrl:    webRepo.WebUrl,
		WebStatus: webRepo.WebStatus,
		IconUrl:   webRepo.IconUrl,
	}
}
