package feature_test

import (
	"github.com/ktsimpso/machine_learning/feature"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyTable(t *testing.T) {
	assert := assert.New(t)
	table := feature.CreateTable([]feature.Feature{})
	assert.Equal(0, table.NumColumns())
	assert.Equal(0, table.NumRows())
}

func TestNumColumns(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)
	assert.Equal(len(featureList), table.NumColumns())
}

func TestLabelFromColumnIndex(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)

	for index, f := range featureList {
		assert.Equal(f.TypeKey(), table.LabelFromColumnIndex(index).TypeKey())
	}
}

func TestColumnIndexFromLabel(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)

	for index, f := range featureList {
		assert.Equal(index, table.ColumnIndexFromLabel(f.TypeKey()))
	}
}

func TestAddStringRow(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)
	row1 := []string{"3.14", "2.0", "456.23", "a", "e", "i"}

	table.AddStringRow(row1)
	assert.Equal(1, table.NumRows())

	for index, item := range row1 {
		assert.Equal(item, table.At(0, index).StringValue)
	}

	row2 := []string{"3.14", "2.0", "456.23", "a", "z", "i"}
	table.AddStringRow(row2)
	assert.Equal(2, table.NumRows())

	for index, item := range row1 {
		if index == 4 {
			assert.Nil(table.At(1, index))
			continue
		}
		assert.Equal(item, table.At(1, index).StringValue)
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail("Expected a panic for an invalid lengthed string row")
			}
		}()
		table.AddStringRow(row2[1:])
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail("Expected a panic for an invalid lengthed string row")
			}
		}()
		table.AddStringRow([]string{})
	}()
}

func TestNumRows(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)
	row := []string{"3.14", "2.0", "456.23", "a", "e", "i"}

	for i := 0; i < 15; i++ {
		assert.Equal(i, table.NumRows())
		table.AddStringRow(row)
	}
}

func TestAt(t *testing.T) {
	assert := assert.New(t)

	table := feature.CreateTable(featureList)
	row1 := []string{"3.14", "2.0", "456.23", "a", "e", "i"}
	row2 := []string{"2.44", "1.2", "111.3", "b", "d", "h"}
	row3 := []string{"31", "42.0", "44444", "c", "f", "g"}
	table.AddStringRow(row1)
	table.AddStringRow(row2)
	table.AddStringRow(row3)

	for index, value := range row1 {
		assert.Equal(value, table.At(0, index).StringValue)
	}

	for index, value := range row2 {
		assert.Equal(value, table.At(1, index).StringValue)
	}

	for index, value := range row3 {
		assert.Equal(value, table.At(2, index).StringValue)
	}
}

func TestAddColumn(t *testing.T) {
	assert := assert.New(t)
	newFeatureIndex := len(featureList) - 1

	table := feature.CreateTable(featureList[:newFeatureIndex])
	row1 := []string{"3.14", "2.0", "456.23", "a", "e"}
	row2 := []string{"2.44", "1.2", "111.3", "b", "d"}
	row3 := []string{"31", "42.0", "44444", "c", "f"}
	table.AddStringRow(row1)
	table.AddStringRow(row2)
	table.AddStringRow(row3)
	assert.Equal(newFeatureIndex, table.NumColumns())

	newFeature := featureList[newFeatureIndex]
	column := []*feature.Instance{
		newFeature.Create("i"),
		newFeature.Create("h"),
		newFeature.Create("g"),
	}

	table.AddColumn(newFeature, column)
	assert.Equal(len(featureList), table.NumColumns())
	assert.Equal(newFeature.TypeKey(), table.LabelFromColumnIndex(newFeatureIndex).TypeKey())
	assert.Equal(newFeatureIndex, table.ColumnIndexFromLabel(newFeature.TypeKey()))

	for index, item := range column {
		assert.Equal(item, table.At(index, newFeatureIndex))
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail("Expected a panic for an invalid lengthed column")
			}
		}()
		table.AddColumn(newFeature, column[:2])
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail("Expected a panic for an invalid lengthed column")
			}
		}()
		table.AddColumn(newFeature, []*feature.Instance{})
	}()
}
