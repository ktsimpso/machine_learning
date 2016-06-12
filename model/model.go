package model

import (
	"errors"
	"github.com/ktsimpso/machine_learning/feature"
	"math"
)

type Model interface {
	Train(data []map[string]feature.Instance, featureList []feature.Feature, resultFeature feature.Feature) error
	Predict(data []map[string]feature.Instance, ft feature.Feature) ([]feature.Instance, error)
}

type DecisionTree struct {
	left          *DecisionTree
	right         *DecisionTree
	switchFeature feature.Instance
	endState      bool
	prediction    string
}

func (dt *DecisionTree) Train(data []map[string]feature.Instance, featureList []feature.Feature, resultFeature feature.Feature) error {
	positiveCount := 0

	for _, row := range data {
		result, ok := row[resultFeature.Name]

		if !ok {
			return errors.New("There was an error when getting the resultFeature")
		}

		//TODO: support more than binary features
		if result.DiscreteValue == 1 {
			positiveCount += 1
		}
	}

	baseScore := float64(positiveCount) / float64(len(data))
	normailzedBaseScore := math.Abs(baseScore - 0.5)
	highScore := 0.0

	// This whole code block makes me die a little...on the inside
	for _, f := range featureList {
		if f.Name == resultFeature.Name {
			continue
		}

		switch f.Type {
		case feature.Discrete:
			type valuesCounters struct {
				positiveCount, totalCount int64
			}
			valuesPositiveCount := map[int64]*valuesCounters{}

			for _, row := range data {
				value, ok := row[f.Name]
				result := row[resultFeature.Name].DiscreteValue
				if !ok {
					continue
				}

				intValue := value.DiscreteValue

				if _, ok := valuesPositiveCount[intValue]; !ok {
					valuesPositiveCount[intValue] = &valuesCounters{0, 0}
				}

				if result == 1 {
					valuesPositiveCount[intValue].positiveCount += 1
				}

				valuesPositiveCount[intValue].totalCount += 1
			}

			for key, valueCount := range valuesPositiveCount {
				score := math.Abs(float64(valueCount.positiveCount)/float64(valueCount.totalCount) - 0.5)
				if score > highScore {
					highScore = score
					sf, err := f.Create(f.ReverseValueMap[key])
					if err != nil {
						panic(err) // Too tired to think about this right now TODO: look into this later
					}
					dt.switchFeature = sf
				}
			}
		case feature.Continuous:
			//TODO: continuous shits
		}
	}

	if highScore > (normailzedBaseScore) {
		leftData := []map[string]feature.Instance{}
		rightData := []map[string]feature.Instance{}

		for _, row := range data {
			value, ok := row[dt.switchFeature.Feature.Name]
			if !ok {
				leftData = append(leftData, row)
				continue
			}

			switch value.Feature.Type {
			case feature.Discrete:
				if value.DiscreteValue == dt.switchFeature.DiscreteValue {
					rightData = append(rightData, row)
				} else {
					leftData = append(leftData, row)
				}
			case feature.Continuous:
				if value.ContinuousValue > dt.switchFeature.ContinuousValue {
					rightData = append(rightData, row)
				} else {
					leftData = append(leftData, row)
				}
			}
		}

		dt.left = &DecisionTree{}
		dt.right = &DecisionTree{}

		dt.left.Train(leftData, featureList, resultFeature)
		dt.right.Train(rightData, featureList, resultFeature)
	} else {
		dt.endState = true
		if baseScore > 0.5 {
			dt.prediction = resultFeature.ReverseValueMap[1]
		} else {
			dt.prediction = resultFeature.ReverseValueMap[0]
		}
	}

	return nil
}

func (dt *DecisionTree) Predict(data []map[string]feature.Instance, ft feature.Feature) ([]feature.Instance, error) {
	results := []feature.Instance{}

	for _, row := range data {
		result, err := dt.predictRow(row, ft)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (dt *DecisionTree) predictRow(row map[string]feature.Instance, ft feature.Feature) (feature.Instance, error) {
	if dt.endState {
		predictionType, err := ft.Create(dt.prediction)
		if err != nil {
			return feature.Instance{}, err
		}
		return predictionType, nil
	}

	determiningFeature, ok := row[dt.switchFeature.Feature.Name]

	if !ok {
		return dt.left.predictRow(row, ft)
	}

	switch determiningFeature.Feature.Type {
	case feature.Discrete:
		if determiningFeature.DiscreteValue == dt.switchFeature.DiscreteValue {
			return dt.right.predictRow(row, ft)
		} else {
			return dt.left.predictRow(row, ft)
		}
	case feature.Continuous:
		if determiningFeature.ContinuousValue > dt.switchFeature.ContinuousValue {
			return dt.right.predictRow(row, ft)
		} else {
			return dt.left.predictRow(row, ft)
		}
	default:
		return feature.Instance{}, errors.New("Unknown Feature Type")
	}
}
