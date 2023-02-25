package logic

import "backup/models"

type Mock struct {
}

func (m Mock) GetTag() string {
	return "Mock"
}

func (m Mock) GetCollector() ([]*models.Collector, error) {
	return []*models.Collector{
		{
			EquipName: "t1",
			Addr:      "192.168.1.2:21",
		},
		{
			EquipName: "t2",
			Addr:      "192.168.1.3:21",
		},
		{
			EquipName: "t3",
			Addr:      "192.168.1.4:21",
		},
		{
			EquipName: "t4",
			Addr:      "192.168.1.5:21",
		},
		{
			EquipName: "t5",
			Addr:      "192.168.1.6:21",
		},
		{
			EquipName: "t6",
			Addr:      "192.168.1.7:21",
		}, {
			EquipName: "t7",
			Addr:      "192.168.1.8:21",
		},
	}, nil
}

var _ QueryCollector = (*Mock)(nil)
