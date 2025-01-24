package device

import (
	"errors"
	"fmt"
)

type List struct {
	devices map[string]Controller
}

func NewController() *List {
	return &List{devices: make(map[string]Controller)}
}

func (list *List) AddDevice(id string, d Controller) Controller {
	if _, exist := list.devices[id]; exist {
		fmt.Printf("Device '%s' already exists\n", id)
		return nil
	}
	list.devices[id] = d
	fmt.Printf("add device: '%s' \n", id)
	return list.devices[id]
}

func (list *List) RemoveDevice(id string) error {
	if _, exist := list.devices[id]; exist {
		delete(list.devices, id)
		fmt.Printf("remove device: '%s' \n", id)
		return nil
	}
	return errors.New("device does not exist")
}
