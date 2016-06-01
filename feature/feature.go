package feature

import (
	"errors"
	"fmt"
	"strconv"
)

type Feature interface {
	Name() string
}

type Type func(value string) (Feature, error)

type baseFeature struct {
	name string
}

func (bf baseFeature) Name() string {
	return bf.name
}

type Continuous struct {
	baseFeature
	Value float64
}

func NewContinousType(name string) Type {
	return func(value string) (Feature, error) {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		return Continuous{
			baseFeature{
				name,
			},
			floatValue,
		}, nil
	}
}

type Discrete struct {
	baseFeature
	Value           int64
	Values          []string
	FeatureType     Type
	ValueMap        map[string]int64
	ReverseValueMap map[int64]string
}

func (d Discrete) String() string {
	return fmt.Sprintf("{name{%s} Value: {%d, %s}}", d.name, d.Value, d.ReverseValueMap[d.Value])
}

func NewDiscreteType(name string, values []string) Type {
	valueMap := map[string]int64{}
	reverseValueMap := map[int64]string{}

	for index, value := range values {
		valueMap[value] = int64(index)
		reverseValueMap[int64(index)] = value
	}

	var featureType Type

	featureType = func(value string) (Feature, error) {
		index, ok := valueMap[value]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Value type: %s not found for DiscreteFeature with name: %s", value, name))
		}

		return Discrete{
			baseFeature{
				name,
			},
			index,
			values,
			featureType,
			valueMap,
			reverseValueMap,
		}, nil
	}

	return featureType
}
