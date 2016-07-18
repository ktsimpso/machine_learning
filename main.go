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

	fmt.Println(data.NumRows())

	dt := model.DecisionTree{}
	err = dt.Train(data, IncomeFeature)
	if err != nil {
		panic(err)
	}

	testData, err := getDataFromFile("data/test.tsv", FeatureList[:len(FeatureList)-1])
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
		output.Guesses = append(output.Guesses, prediction.StringValue)
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

func getDataFromFile(filename string, features []feature.Feature) (*feature.Table, error) {
	dataFile, err := os.Open(filename)
	if err != nil {
		return &feature.Table{}, err
	}
	defer dataFile.Close()

	csvReader := csv.NewReader(dataFile)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = len(features)

	data := feature.CreateTable(features)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &feature.Table{}, err
		}

		data.AddStringRow(row)
	}

	for _, f := range FeatureList {
		if f.Type == feature.Discrete {
			continue
		}
		data.AddColumn(feature.ConvertContinuousToDiscrete(data.GetColumn(f.TypeKey())))
	}

	return data, nil
}
