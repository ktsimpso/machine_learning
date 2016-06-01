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
	assert.Nil(err)

	trainingData := []map[string]feature.Feature{}
	testData := []map[string]feature.Feature{}

	for _, datum := range data {
		if rand.Float64() < 0.5 {
			trainingData = append(trainingData, datum)
		} else {
			testData = append(testData, datum)
		}
	}

	dt := model.DecisionTree{}
	err = dt.Train(trainingData, ConcreteFeatureList, ConcreteFeatureList[14].(feature.Discrete))
	assert.Nil(err)

	predictions, err := dt.Predict(testData, IncomeFeature)
	assert.Nil(err)

	count := 0
	correct := 0

	for index, testDatum := range testData {
		if (predictions[index]).Value == (testDatum["income"]).(feature.Discrete).Value {
			correct += 1
		}

		count += 1
	}

	result := float64(correct) / float64(count)
	assert.True(result > 0.75)
	fmt.Println(result)
}
