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
	err = dt.Train(data, ConcreteFeatureList, ConcreteFeatureList[14].(feature.Discrete))
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
		output.Guesses = append(output.Guesses, prediction.ReverseValueMap[prediction.Value])
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

func getDataFromFile(filename string, features []feature.Type) ([]map[string]feature.Feature, error) {
	dataFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()

	csvReader := csv.NewReader(dataFile)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = len(features)

	data := []map[string]feature.Feature{}

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

func recordToFeatures(record []string, featureList []feature.Type) map[string]feature.Feature {
	features := map[string]feature.Feature{}

	for index, featureType := range featureList {
		f, err := featureType(record[index])
		if err != nil {
			//fmt.Println(err)
			continue
		}

		features[f.Name()] = f
	}

	return features
}
