package model

import (
	"github.com/ktsimpso/machine_learning/feature"
	"math"
)

type Model interface {
	Train(data feature.TableViewer, resultFeature feature.Feature) error
	Predict(data feature.TableViewer, resultFeature feature.Feature) ([]feature.Instance, error)
}

type DecisionTree struct {
	left          *DecisionTree
	right         *DecisionTree
	switchFeature *feature.Instance
	endState      bool
	prediction    *feature.Instance
}

func (dt *DecisionTree) Train(data feature.TableViewer, resultFeature feature.Feature) error {
	positiveCount := 0
	resultColumnIndex := data.ColumnIndexFromLabel(resultFeature.TypeKey())

	for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
		//TODO: support more than brinary features
		if data.At(rowIndex, resultColumnIndex).DiscreteValue == 1 {
			positiveCount += 1
		}
	}

	baseScore := float64(positiveCount) / float64(data.NumRows())
	normailzedBaseScore := math.Abs(baseScore - 0.5)
	highScore := 0.0

	for columnIndex := 0; columnIndex < data.NumColumns(); columnIndex++ {
		if columnIndex == resultColumnIndex {
			continue
		}

		switch data.LabelFromColumnIndex(columnIndex).Type {
		case feature.Discrete:
			type valuesCounters struct {
				positiveCount, totalCount int64
			}
			valuesPositiveCount := map[int64]*valuesCounters{}

			for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
				instance := data.At(rowIndex, columnIndex)
				if instance == nil {
					continue
				}

				instanceValue := instance.DiscreteValue

				if _, ok := valuesPositiveCount[instanceValue]; !ok {
					valuesPositiveCount[instanceValue] = &valuesCounters{0, 0}
				}

				if data.At(rowIndex, resultColumnIndex).DiscreteValue == 1 {
					valuesPositiveCount[instanceValue].positiveCount += 1
				}

				valuesPositiveCount[instanceValue].totalCount += 1
			}

			for key, valueCount := range valuesPositiveCount {
				score := math.Abs(float64(valueCount.positiveCount)/float64(valueCount.totalCount) - 0.5)
				if score > highScore {
					highScore = score
					sf := data.LabelFromColumnIndex(columnIndex).CreateDiscrete(key)
					dt.switchFeature = sf
				}
			}
		}
	}

	if highScore > normailzedBaseScore {
		leftData := feature.NewTableViewBuilder(data).WithAllColumns()
		rightData := feature.NewTableViewBuilder(data).WithAllColumns()
		switchColumnIndex := data.ColumnIndexFromLabel(dt.switchFeature.Feature.TypeKey())

		for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
			if isFeatureRight(data.At(rowIndex, switchColumnIndex), dt.switchFeature) {
				rightData.WithRow(rowIndex)
			} else {
				leftData.WithRow(rowIndex)
			}
		}

		dt.left = &DecisionTree{}
		dt.right = &DecisionTree{}

		dt.left.Train(leftData.Build(), resultFeature)
		dt.right.Train(rightData.Build(), resultFeature)
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

func (dt *DecisionTree) Predict(data feature.TableViewer, resultFeature feature.Feature) ([]feature.Instance, error) {
	results := []feature.Instance{}

	for rowIndex := 0; rowIndex < data.NumRows(); rowIndex++ {
		result, err := dt.predictRow(data, rowIndex, resultFeature)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (dt *DecisionTree) predictRow(data feature.TableViewer, rowIndex int, resultFeature feature.Feature) (feature.Instance, error) {
	if dt.endState {
		return *dt.prediction, nil
	}

	columnIndex := data.ColumnIndexFromLabel(dt.switchFeature.Feature.TypeKey())

	if isFeatureRight(data.At(rowIndex, columnIndex), dt.switchFeature) {
		return dt.right.predictRow(data, rowIndex, resultFeature)
	} else {
		return dt.left.predictRow(data, rowIndex, resultFeature)
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
