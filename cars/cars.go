package cars

import (
	"encoding/json"
	"os"
)

var DefaultCar = Car{
	Maker:      "Unknown",
	Model:      "Unknown",
	Group:      "Unknown",
	Year:       9999,
	CarOrdinal: 0,
	Weight:     9999,
}

type Car struct {
	CarOrdinal int32 `json:"car_ordinal"`
	Group      string `json:"group"`
	Maker      string `json:"maker"`
	Model      string `json:"model"`
	Year       int32 `json:"year"`
	Weight     int32 `json:"weight"`
}

func FindCar(a []Car, x int32) int32 {
	for i, n := range a {
		if x == n.CarOrdinal {
			return int32(i)
		}
	}
	return -1
}
func HasCarChanged(old int32,new int32) bool {
  return new == old

}

func SetCurrentCar(cars []Car,id int32) Car {
  car := FindCar(cars, id)
				if car == -1 {
					return DefaultCar
				} else {
					return cars[car]
				}
}

func ReadCarList(path string) ([]Car, error) {
	cars := make([]Car,700)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
  
	err = json.Unmarshal(content, &cars)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

