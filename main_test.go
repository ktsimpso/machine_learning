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
	fmt.Println(len(data.Rows))
	assert.Nil(err)

	trainingData := feature.CreateTable(data.Features)
	testData := feature.CreateTable(data.Features)

	for _, row := range data.Rows {
		if rand.Float64() < 0.5 {
			trainingData.AddRow(row)
		} else {
			testData.AddRow(row)
		}
	}

	dt := model.DecisionTree{}
	err = dt.Train(trainingData, IncomeFeature)
	assert.Nil(err)

	predictions, err := dt.Predict(testData, IncomeFeature)
	assert.Nil(err)

	count := 0
	correct := 0

	for index, testDatum := range testData.Columns[testData.FeatureMap[IncomeFeature.TypeKey()]] {
		if (predictions[index]).DiscreteValue == testDatum.DiscreteValue {
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

	for index, trainingDatum := range trainingData.Columns[trainingData.FeatureMap[IncomeFeature.TypeKey()]] {
		if (predictions[index]).DiscreteValue == trainingDatum.DiscreteValue {
			correct += 1
		}

		count += 1
	}

	result = float64(correct) / float64(count)
	assert.True(result > 0.75)
	fmt.Println(result)
}
