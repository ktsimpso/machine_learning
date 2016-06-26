package model

import (
	"github.com/ktsimpso/machine_learning/feature"
	"math"
)

type Model interface {
	Train(data *feature.Table, resultFeature feature.Feature) error
	Predict(data *feature.Table, resultFeature feature.Feature) ([]feature.Instance, error)
}

type DecisionTree struct {
	left          *DecisionTree
	right         *DecisionTree
	switchFeature *feature.Instance
	endState      bool
	prediction    *feature.Instance
}

func (dt *DecisionTree) Train(data *feature.Table, resultFeature feature.Feature) error {
	positiveCount := 0
	resultColumn := data.Columns[data.FeatureMap[resultFeature.TypeKey()]]

	for _, result := range resultColumn {
		//TODO: support more than binary features
		if result.DiscreteValue == 1 {
			positiveCount += 1
		}
	}

	baseScore := float64(positiveCount) / float64(len(data.Rows))
	normailzedBaseScore := math.Abs(baseScore - 0.5)
	highScore := 0.0

	for _, f := range data.Features {
		if f.Name == resultFeature.Name {
			continue
		}

		switch f.Type {
		case feature.Discrete:
			type valuesCounters struct {
				positiveCount, totalCount int64
			}
			valuesPositiveCount := map[int64]*valuesCounters{}

			for index, instance := range data.Columns[data.FeatureMap[f.TypeKey()]] {
				if instance == nil {
					continue
				}

				instanceValue := instance.DiscreteValue
				resultValue := resultColumn[index].DiscreteValue

				if _, ok := valuesPositiveCount[instanceValue]; !ok {
					valuesPositiveCount[instanceValue] = &valuesCounters{0, 0}
				}

				if resultValue == 1 {
					valuesPositiveCount[instanceValue].positiveCount += 1
				}

				valuesPositiveCount[instanceValue].totalCount += 1
			}

			for key, valueCount := range valuesPositiveCount {
				score := math.Abs(float64(valueCount.positiveCount)/float64(valueCount.totalCount) - 0.5)
				if score > highScore {
					highScore = score
					sf := f.Create(f.ReverseValueMap[key])
					dt.switchFeature = sf
				}
			}
		case feature.Continuous:
			//TODO: this method of continuous functions doesn't seem to perform at all! Should probably come up with something better
			max := -math.MaxFloat64
			min := math.MaxFloat64

			continuousColumn := data.Columns[data.FeatureMap[f.TypeKey()]]

			for _, instance := range continuousColumn {
				if instance == nil {
					continue
				}
				max = math.Max(max, instance.ContinuousValue)
				min = math.Min(min, instance.ContinuousValue)
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

				for index, instance := range continuousColumn {
					if instance == nil {
						continue
					}

					resultValue := resultColumn[index].DiscreteValue

					if instance.ContinuousValue < pivot {
						leftTotal += 1
						if resultValue == 1 {
							leftPositiveCount += 1
						}
					} else {
						rightTotal += 1
						if resultValue == 1 {
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
				dt.switchFeature = f.CreateContinuous(lastPivot)
			}
		}
	}

	if highScore > normailzedBaseScore {
		leftData := feature.CreateTable(data.Features)
		rightData := feature.CreateTable(data.Features)

		for index, instance := range data.Columns[data.FeatureMap[dt.switchFeature.Feature.TypeKey()]] {
			var appendData *feature.Table

			if isFeatureRight(instance, dt.switchFeature) {
				appendData = rightData
			} else {
				appendData = leftData
			}

			appendData.AddRow(data.Rows[index])
		}

		dt.left = &DecisionTree{}
		dt.right = &DecisionTree{}

		dt.left.Train(leftData, resultFeature)
		dt.right.Train(rightData, resultFeature)
	} else {
		dt.endState = true
		var err error

		if baseScore > 0.5 {
			dt.prediction = resultFeature.Create(resultFeature.ReverseValueMap[1])
		} else {
			dt.prediction = resultFeature.Create(resultFeature.ReverseValueMap[0])
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (dt *DecisionTree) Predict(data *feature.Table, resultFeature feature.Feature) ([]feature.Instance, error) {
	results := []feature.Instance{}

	for _, row := range data.Rows {
		result, err := dt.predictRow(data, row, resultFeature)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (dt *DecisionTree) predictRow(data *feature.Table, row []*feature.Instance, resultFeature feature.Feature) (feature.Instance, error) {
	if dt.endState {
		return *dt.prediction, nil
	}

	determiningFeature := row[data.FeatureMap[dt.switchFeature.Feature.TypeKey()]]

	if isFeatureRight(determiningFeature, dt.switchFeature) {
		return dt.right.predictRow(data, row, resultFeature)
	} else {
		return dt.left.predictRow(data, row, resultFeature)
	}
}

func isFeatureRight(testFeature, switchFeature *feature.Instance) bool {
	if testFeature == nil {
		return false
	}

	switch testFeature.Feature.Type {
	case feature.Discrete:
		return testFeature.DiscreteValue == switchFeature.DiscreteValue
	case feature.Continuous:
		return testFeature.ContinuousValue > switchFeature.ContinuousValue
	default:
		return false
	}
}
