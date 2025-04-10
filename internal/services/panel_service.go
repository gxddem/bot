package services

import "foxhole-bot/internal/entities"



type PanelStorage interface{
	CreateList(string) (*entities.Panel, error)
	GetList(string) (*entities.Panel, error)
	RemoveItemFromList(*entities.Panel, int) error
}

type PanelService struct {
	Storage PanelStorage
}

func NewPanelService(storage PanelStorage) *PanelService {
	return &PanelService{
		Storage: storage,
	}
}


func (ps *PanelService) CreateList (channelID string) (*entities.Panel, error) {
	panel, err := ps.Storage.CreateList(channelID)
	if err != nil {
		return nil, err
	}
	
	return panel, nil
}

func (ps *PanelService) GetList (channelID string) (*entities.Panel, error) {
	panel, err := ps.Storage.GetList(channelID)
	if err != nil {
		return nil, err
	}

	return panel, nil
}

func (ps *PanelService) AppendItemToList (panel *entities.Panel, data string) error {
	panel.Items = append(panel.Items, data)
	return nil 
}

func (ps *PanelService) RemoveItemFromList (panel *entities.Panel, index int) error {
	return ps.Storage.RemoveItemFromList(panel, index)
}