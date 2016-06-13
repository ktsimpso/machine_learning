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
	prediction    feature.Instance
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
						panic(err)
					}
					dt.switchFeature = sf
				}
			}
		case feature.Continuous:
			//TODO: this method of continuous functions doesn't seem to perform at all! Should probably come up with something better
			max := -math.MaxFloat64
			min := math.MaxFloat64

			for _, row := range data {
				value, ok := row[f.Name]
				if !ok {
					continue
				}

				floatValue := value.ContinuousValue
				max = math.Max(max, floatValue)
				min = math.Min(min, floatValue)
			}

			pivot := (max + min) / 2
			lastPivot := pivot
			bestScore := 0.0
			lastBestScore := 0.0

			for {
				leftPositiveCount := 0
				leftTotal := 0
				rightPositiveCount := 0
				rightTotal := 0

				for _, row := range data {
					value, ok := row[f.Name]
					result := row[resultFeature.Name].DiscreteValue
					if !ok {
						continue
					}

					floatValue := value.ContinuousValue

					if floatValue < pivot {
						leftTotal += 1
						if result == 1 {
							leftPositiveCount += 1
						}
					} else {
						rightTotal += 1
						if result == 1 {
							rightPositiveCount += 1
						}
					}
				}

				leftScore := math.Abs(float64(leftPositiveCount)/float64(leftTotal) - 0.5)
				rightScore := math.Abs(float64(rightPositiveCount)/float64(rightTotal) - 0.5)
				bestScore = math.Max(leftScore, rightScore)

				if bestScore > lastBestScore {
					lastPivot = pivot
					lastBestScore = bestScore

					if leftScore > rightScore {
						pivot = (min + pivot) / 2
					} else {
						pivot = (max + pivot) / 2
					}
				} else {
					break
				}
			}

			if bestScore > highScore {
				highScore = bestScore
				sf, err := f.CreateContinuous(lastPivot)
				if err != nil {
					panic(err)
				}
				dt.switchFeature = sf
			}

		}
	}

	if highScore > normailzedBaseScore {
		leftData := []map[string]feature.Instance{}
		rightData := []map[string]feature.Instance{}

		for _, row := range data {
			value, ok := row[dt.switchFeature.Feature.Name]
			if !ok {
				leftData = append(leftData, row)
				continue
			}

			if isFeatureRight(value, dt.switchFeature) {
				rightData = append(rightData, row)
			} else {
				leftData = append(leftData, row)
			}
		}

		dt.left = &DecisionTree{}
		dt.right = &DecisionTree{}

		dt.left.Train(leftData, featureList, resultFeature)
		dt.right.Train(rightData, featureList, resultFeature)
	} else {
		dt.endState = true
		var err error

		if baseScore > 0.5 {
			dt.prediction, err = resultFeature.Create(resultFeature.ReverseValueMap[1])
		} else {
			dt.prediction, err = resultFeature.Create(resultFeature.ReverseValueMap[0])
		}

		if err != nil {
			return err
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
		return dt.prediction, nil
	}

	determiningFeature, ok := row[dt.switchFeature.Feature.Name]

	if !ok {
		return dt.left.predictRow(row, ft)
	}

	if isFeatureRight(determiningFeature, dt.switchFeature) {
		return dt.right.predictRow(row, ft)
	} else {
		return dt.left.predictRow(row, ft)
	}
}

func isFeatureRight(testFeature, switchFeature feature.Instance) bool {
	switch testFeature.Feature.Type {
	case feature.Discrete:
		return testFeature.DiscreteValue == switchFeature.DiscreteValue
	case feature.Continuous:
		return testFeature.ContinuousValue > switchFeature.ContinuousValue
	default:
		return false
	}
}
