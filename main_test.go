package main

import (
	"fmt"
	"github.com/ktsimpso/machine_learning/feature"
	"github.com/ktsimpso/machine_learning/model"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UTC().UnixNano())

	data, err := getDataFromFile("data/training.tsv", FeatureList)
	fmt.Println(data.NumRows())
	assert.Nil(err)

	trainingDataBuilder := feature.NewTableViewBuilder(data).WithAllColumns()
	testDataBuidler := feature.NewTableViewBuilder(data).WithAllColumns()

	for row := range data.Rows() {
		if rand.Float64() < 0.5 {
			trainingDataBuilder.WithRow(row.Index())
		} else {
			testDataBuidler.WithRow(row.Index())
		}
	}

	trainingData := trainingDataBuilder.Build()
	testData := testDataBuidler.Build()

	dt := model.DecisionTree{}
	err = dt.Train(trainingData, IncomeFeature)
	assert.Nil(err)

	predictions, err := dt.Predict(testData, IncomeFeature)
	assert.Nil(err)

	count := 0
	correct := 0

	for record := range testData.GetColumn(IncomeFeature.TypeKey()).Instances() {
		if (predictions[record.Index]).DiscreteValue == record.Instance.DiscreteValue {
			correct += 1
		}

		count += 1
	}

	result := float64(correct) / float64(count)
	assert.True(result > 0.75)
	fmt.Println(result)

	predictions, err = dt.Predict(trainingData, IncomeFeature)
	assert.Nil(err)
	count = 0
	correct = 0

	for record := range trainingData.GetColumn(IncomeFeature.TypeKey()).Instances() {
		if (predictions[record.Index]).DiscreteValue == record.Instance.DiscreteValue {
			correct += 1
		}

		count += 1
	}

	result = float64(correct) / float64(count)
	assert.True(result > 0.75)
	fmt.Println(result)
}
