package feature

import (
	"math"
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
	CreateDiscrete   CreateDiscrete
}

type Instance struct {
	Feature         Feature
	DiscreteValue   int64
	ContinuousValue float64
	StringValue     string
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
type CreateDiscrete func(value int64) *Instance

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
				value,
			}
		},
		func(value float64) *Instance {
			return &Instance{
				this,
				0,
				value,
				strconv.FormatFloat(value, 'f', -1, 64),
			}
		},
		func(value int64) *Instance {
			return nil
		},
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
				value,
			}
		},
		func(value float64) *Instance {
			return nil
		},
		func(value int64) *Instance {
			stringValue, ok := reverseValueMap[value]
			if !ok {
				return nil
			}

			return &Instance{
				this,
				value,
				0.0,
				stringValue,
			}
		},
	}

	return this
}

func ConvertContinuousToDiscrete(columnIndex int, data *Table) (Feature, []*Instance) {
	total := 0.0
	count := 0
	columnFeature := data.LabelFromColumnIndex(columnIndex)

	//TODO: validation
	for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
		record := data.At(rowIndex, columnIndex)
		if record == nil {
			continue
		}

		//TODO: this could overflow
		total += record.ContinuousValue
		count += 1
	}

	mean := total / float64(count)
	squaredDistance := 0.0

	for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
		record := data.At(rowIndex, columnIndex)
		if record == nil {
			continue
		}

		distance := record.ContinuousValue - mean
		squaredDistance += distance * distance
	}

	standardDeviation := math.Sqrt(squaredDistance / (total - 1))

	var this Feature
	this = Feature{
		Discrete,
		columnFeature.Name,
		func(value string) *Instance {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil
			}

			intValue := int64((floatValue - mean) / standardDeviation)

			return &Instance{
				this,
				intValue,
				floatValue,
				value,
			}
		},
		func(value float64) *Instance {
			intValue := int64((value - mean) / standardDeviation)

			return &Instance{
				this,
				intValue,
				value,
				strconv.FormatFloat(value, 'f', -1, 64),
			}
		},
		func(value int64) *Instance {
			floatValue := mean + float64(value)*standardDeviation

			if value >= 0 {
				floatValue += (standardDeviation / 2)
			} else {
				floatValue -= (standardDeviation / 2)
			}

			return &Instance{
				this,
				value,
				floatValue,
				strconv.FormatFloat(floatValue, 'f', -1, 64),
			}
		},
	}

	convertedFeatures := make([]*Instance, data.NumRows())

	for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
		record := data.At(rowIndex, columnIndex)
		if record == nil {
			convertedFeatures[rowIndex] = nil
			continue
		}

		convertedFeatures[rowIndex] = this.CreateContinuous(record.ContinuousValue)
	}

	return this, convertedFeatures
}
