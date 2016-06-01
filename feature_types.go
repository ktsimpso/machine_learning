package main

import (
	"github.com/ktsimpso/machine_learning/feature"
)

var FeatureList = []feature.Type{
	AgeFeature,
	WorkClassFeature,
	SampleWeightFeature,
	EducationFeature,
	EducationNumberFeature,
	MaritalStatusFeature,
	OccupationFeature,
	RelationshipFeature,
	RaceFeature,
	SexFeature,
	CapitalGainFeature,
	CapitalLossFeature,
	HoursPerWeekFeature,
	NativeCountryFeature,
	IncomeFeature,
}

//TODO: remove this once I fix the types....
var ConcreteFeatureList = []feature.Feature{}

func init() {
	f, err := AgeFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = WorkClassFeature("Private")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = SampleWeightFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = EducationFeature("Bachelors")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = EducationNumberFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = MaritalStatusFeature("Divorced")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = OccupationFeature("Sales")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = RelationshipFeature("Wife")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = RaceFeature("Other")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = SexFeature("Female")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = CapitalGainFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = CapitalLossFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = HoursPerWeekFeature("1")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = NativeCountryFeature("Cambodia")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)

	f, err = IncomeFeature(">50K")
	if err != nil {
		panic(err)
	}
	ConcreteFeatureList = append(ConcreteFeatureList, f)
}

//End TODO

var AgeFeature = feature.NewContinousType("age")
var SampleWeightFeature = feature.NewContinousType("sampleWeight")
var EducationNumberFeature = feature.NewContinousType("educationNumber")
var CapitalGainFeature = feature.NewContinousType("capitalGain")
var CapitalLossFeature = feature.NewContinousType("capitalLoss")
var HoursPerWeekFeature = feature.NewContinousType("HoursPerWeek")

var WorkClassFeature = feature.NewDiscreteType("workClass", []string{
	"Private",
	"Self-emp-not-inc",
	"Self-emp-inc",
	"Federal-gov",
	"Local-gov",
	"State-gov",
	"Without-pay",
	"Never-worked",
})

var EducationFeature = feature.NewDiscreteType("education", []string{
	"Bachelors",
	"Some-college",
	"11th",
	"HS-grad",
	"Prof-school",
	"Assoc-acdm",
	"Assoc-voc",
	"9th",
	"7th-8th",
	"12th",
	"Masters",
	"1st-4th",
	"10th",
	"Doctorate",
	"5th-6th",
	"Preschool",
})

var MaritalStatusFeature = feature.NewDiscreteType("maritalStatus", []string{
	"Married-civ-spouse",
	"Divorced",
	"Never-married",
	"Separated",
	"Widowed",
	"Married-spouse-absent",
	"Married-AF-spouse",
})

var OccupationFeature = feature.NewDiscreteType("occupation", []string{
	"Tech-support",
	"Craft-repair",
	"Other-service",
	"Sales",
	"Exec-managerial",
	"Prof-specialty",
	"Handlers-cleaners",
	"Machine-op-inspct",
	"Adm-clerical",
	"Farming-fishing",
	"Transport-moving",
	"Priv-house-serv",
	"Protective-serv",
	"Armed-Forces",
})

var RelationshipFeature = feature.NewDiscreteType("relationship", []string{
	"Wife",
	"Own-child",
	"Husband",
	"Not-in-family",
	"Other-relative",
	"Unmarried",
})

var RaceFeature = feature.NewDiscreteType("race", []string{
	"White",
	"Asian-Pac-Islander",
	"Amer-Indian-Eskimo",
	"Other",
	"Black",
})

var SexFeature = feature.NewDiscreteType("sex", []string{
	"Female",
	"Male",
})

var NativeCountryFeature = feature.NewDiscreteType("nativeCountry", []string{
	"United-States",
	"Cambodia",
	"England",
	"Puerto-Rico",
	"Canada",
	"Germany",
	"Outlying-US(Guam-USVI-etc)",
	"India",
	"Japan",
	"Greece",
	"South",
	"China",
	"Cuba",
	"Iran",
	"Honduras",
	"Philippines",
	"Italy",
	"Poland",
	"Jamaica",
	"Vietnam",
	"Mexico",
	"Portugal",
	"Ireland",
	"France",
	"Dominican-Republic",
	"Laos",
	"Ecuador",
	"Taiwan",
	"Haiti",
	"Columbia",
	"Hungary",
	"Guatemala",
	"Nicaragua",
	"Scotland",
	"Thailand",
	"Yugoslavia",
	"El-Salvador",
	"Trinadad&Tobago",
	"Peru",
	"Hong",
	"Holand-Netherlands",
})

var IncomeFeature = feature.NewDiscreteType("income", []string{
	">50K",
	"<=50K",
})
