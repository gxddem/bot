package storage

import (
	"errors"
	"foxhole-bot/internal/entities"
)

type PanelStorage struct {
	Panels map[string]*entities.Panel
}

func NewStorage() *PanelStorage {
	return &PanelStorage{
		Panels: make(map[string]*entities.Panel),
	}
}

func (ps *PanelStorage) CreateList(channelID string) (*entities.Panel, error) {
	_, exists := ps.Panels[channelID]
	if exists {
		return nil, errors.New("a list for this channel already exists")
	}
	ps.Panels[channelID] = &entities.Panel{ChannelID: channelID, Items: make([]string, 0)}

	return ps.Panels[channelID], nil
}

func (ps *PanelStorage) GetList(channelID string) (*entities.Panel, error) {
	panel, exists := ps.Panels[channelID]
	if !exists {
		return nil, errors.New("no list for this channel exists")
	}

	return panel, nil
}

func (ps *PanelStorage) RemoveItemFromList(panel *entities.Panel, index int) error {
	zeroBasedIndex := index - 1
	if zeroBasedIndex < 0 || zeroBasedIndex >= len(panel.Items) {
		return errors.New("no such item")
	}

	panel.Items = append(panel.Items[:zeroBasedIndex], panel.Items[zeroBasedIndex+1:]...)

	return nil
}

func (ps *PanelStorage) PurgeList(panel *entities.Panel) error {
	panel.Items = make([]string, 0)

	return nil
}