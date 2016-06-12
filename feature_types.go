package main

import (
	"github.com/ktsimpso/machine_learning/feature"
)

var FeatureList = []feature.Feature{
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

var AgeFeature = feature.NewContinous("age")
var SampleWeightFeature = feature.NewContinous("sampleWeight")
var EducationNumberFeature = feature.NewContinous("educationNumber")
var CapitalGainFeature = feature.NewContinous("capitalGain")
var CapitalLossFeature = feature.NewContinous("capitalLoss")
var HoursPerWeekFeature = feature.NewContinous("HoursPerWeek")

var WorkClassFeature = feature.NewDiscrete("workClass", []string{
	"Private",
	"Self-emp-not-inc",
	"Self-emp-inc",
	"Federal-gov",
	"Local-gov",
	"State-gov",
	"Without-pay",
	"Never-worked",
})

var EducationFeature = feature.NewDiscrete("education", []string{
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

var MaritalStatusFeature = feature.NewDiscrete("maritalStatus", []string{
	"Married-civ-spouse",
	"Divorced",
	"Never-married",
	"Separated",
	"Widowed",
	"Married-spouse-absent",
	"Married-AF-spouse",
})

var OccupationFeature = feature.NewDiscrete("occupation", []string{
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

var RelationshipFeature = feature.NewDiscrete("relationship", []string{
	"Wife",
	"Own-child",
	"Husband",
	"Not-in-family",
	"Other-relative",
	"Unmarried",
})

var RaceFeature = feature.NewDiscrete("race", []string{
	"White",
	"Asian-Pac-Islander",
	"Amer-Indian-Eskimo",
	"Other",
	"Black",
})

var SexFeature = feature.NewDiscrete("sex", []string{
	"Female",
	"Male",
})

var NativeCountryFeature = feature.NewDiscrete("nativeCountry", []string{
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

var IncomeFeature = feature.NewDiscrete("income", []string{
	">50K",
	"<=50K",
})
