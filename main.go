package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ktsimpso/machine_learning/feature"
	"github.com/ktsimpso/machine_learning/model"
	"io"
	"io/ioutil"
	"os"
)

type Output struct {
	Guesses []string `json:"guesses"`
}

func main() {
	data, err := getDataFromFile("data/training.tsv", FeatureList)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(data))

	dt := model.DecisionTree{}
	err = dt.Train(data, FeatureList, IncomeFeature)
	if err != nil {
		panic(err)
	}

	testData, err := getDataFromFile("data/test.tsv", FeatureList[:14])
	if err != nil {
		panic(err)
	}

	predictions, err := dt.Predict(testData, IncomeFeature)
	if err != nil {
		panic(err)
	}

	output := Output{
		Guesses: []string{},
	}

	for _, prediction := range predictions {
		output.Guesses = append(output.Guesses, prediction.Feature.ReverseValueMap[prediction.DiscreteValue])
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("data/guesses.json", jsonOutput, 0644)
	if err != nil {
		panic(err)
	}
}

func getDataFromFile(filename string, features []feature.Feature) ([]map[string]feature.Instance, error) {
	dataFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()

	csvReader := csv.NewReader(dataFile)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = len(features)

	data := []map[string]feature.Instance{}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data = append(data, recordToFeatures(record, features))
	}

	return data, nil
}

func recordToFeatures(record []string, featureList []feature.Feature) map[string]feature.Instance {
	features := map[string]feature.Instance{}

	for index, feature := range featureList {
		f, err := feature.Create(record[index])
		if err != nil {
			//fmt.Println(err)
			continue
		}

		features[feature.Name] = f
	}

	return features
}
