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
					sf := f.CreateDiscrete(key)
					dt.switchFeature = sf
				}
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
			dt.prediction = resultFeature.CreateDiscrete(1)
		} else {
			dt.prediction = resultFeature.CreateDiscrete(0)
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
