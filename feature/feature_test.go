package feature_test

import (
	"github.com/ktsimpso/machine_learning/feature"
	//"github.com/stretchr/testify/assert"
	//"testing"
)

var featureList = []feature.Feature{
	continuousFeature1,
	continuousFeature2,
	continuousFeature3,
	discreteFeature1,
	discreteFeature2,
	discreteFeature3,
}

var continuousFeature1 = feature.NewContinous("continuousFeature1")
var continuousFeature2 = feature.NewContinous("continuousFeature2")
var continuousFeature3 = feature.NewContinous("continuousFeature3")
var discreteFeature1 = feature.NewDiscrete("discreteFeature1", []string{
	"a",
	"b",
	"c",
})
var discreteFeature2 = feature.NewDiscrete("discreteFeature2", []string{
	"d",
	"e",
	"f",
})
var discreteFeature3 = feature.NewDiscrete("discreteFeature3", []string{
	"g",
	"h",
	"i",
})
