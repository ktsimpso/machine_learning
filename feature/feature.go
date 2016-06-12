package feature

import (
	"errors"
	"fmt"
	"strconv"
)

type FeatureType int

const (
	Continuous FeatureType = iota
	Discrete
)

type Feature struct {
	Type            FeatureType
	Name            string
	Create          Create
	ValueMap        map[string]int64
	ReverseValueMap map[int64]string
}

type Instance struct {
	Feature         Feature
	DiscreteValue   int64
	ContinuousValue float64
}

type Create func(value string) (Instance, error)

func NewContinous(name string) Feature {
	var this Feature
	this = Feature{
		Continuous,
		name,
		func(value string) (Instance, error) {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return Instance{}, err
			}

			return Instance{
				this,
				0,
				floatValue,
			}, nil
		},
		map[string]int64{},
		map[int64]string{},
	}
	return this
}

func NewDiscrete(name string, values []string) Feature {
	valueMap := map[string]int64{}
	reverseValueMap := map[int64]string{}

	for index, value := range values {
		valueMap[value] = int64(index)
		reverseValueMap[int64(index)] = value
	}

	var this Feature
	this = Feature{
		Discrete,
		name,
		func(value string) (Instance, error) {
			index, ok := valueMap[value]
			if !ok {
				return Instance{}, errors.New(fmt.Sprintf("Value type: %s not found for Discrete with name: %s", value, name))
			}

			return Instance{
				this,
				index,
				0.0,
			}, nil
		},
		valueMap,
		reverseValueMap,
	}

	return this
}
