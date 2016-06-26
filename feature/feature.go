package feature

import (
	"strconv"
)

type FeatureType int

const (
	Continuous FeatureType = iota
	Discrete
)

type Feature struct {
	Type             FeatureType
	Name             string
	Create           Create
	CreateContinuous CreateContinuous
	ValueMap         map[string]int64
	ReverseValueMap  map[int64]string
}

type Instance struct {
	Feature         Feature
	DiscreteValue   int64
	ContinuousValue float64
}

type TypeKey struct {
	Name string
	Type FeatureType
}

func (f Feature) TypeKey() TypeKey {
	return TypeKey{
		f.Name,
		f.Type,
	}
}

type Create func(value string) *Instance
type CreateContinuous func(value float64) *Instance

func NewContinous(name string) Feature {
	var this Feature
	this = Feature{
		Continuous,
		name,
		func(value string) *Instance {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil
			}

			return &Instance{
				this,
				0,
				floatValue,
			}
		},
		func(value float64) *Instance {
			return &Instance{
				this,
				0,
				value,
			}
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
		func(value string) *Instance {
			index, ok := valueMap[value]
			if !ok {
				return nil
			}

			return &Instance{
				this,
				index,
				0.0,
			}
		},
		func(value float64) *Instance {
			return nil
		},
		valueMap,
		reverseValueMap,
	}

	return this
}
