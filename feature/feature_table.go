package feature

import (
	"fmt"
)

type Table struct {
	features   []Feature
	featureMap map[TypeKey]int
	items      []*Instance
	numColumns int
	numRows    int
}

func CreateTable(features []Feature) *Table {
	table := Table{}
	featureMap := map[TypeKey]int{}

	for index, feature := range features {
		featureMap[feature.TypeKey()] = index
	}

	table.features = features
	table.featureMap = featureMap
	table.items = []*Instance{}
	table.numColumns = len(features)

	return &table
}

func (t *Table) AddStringRow(records []string) {
	if len(records) != t.numColumns {
		panic(fmt.Sprintf("Number of records added does not equal the number of columns. Expected: %d; Got: %d", t.numColumns, len(records)))
	}

	row := make([]*Instance, len(t.features))

	for index, feature := range t.features {
		row[index] = feature.Create(records[index])
	}

	t.items = append(t.items, row...)
	t.numRows += 1
}

func (t *Table) AddColumn(feature Feature, column []*Instance) {
	if len(column) != t.NumRows() {
		panic(fmt.Sprintf("Number of rows in added column does not equal the number of rows. Expected: %d; Got: %s", t.NumRows(), len(column)))
	}

	t.features = append(t.features, feature)
	t.featureMap[feature.TypeKey()] = len(t.features) - 1

	items := make([]*Instance, len(t.items)+len(column))

	for rowIndex := 0; rowIndex < t.NumRows(); rowIndex++ {
		copy(items[rowIndex*t.numColumns+rowIndex:rowIndex*t.numColumns+t.numColumns+rowIndex], t.items[rowIndex*t.numColumns:rowIndex*t.numColumns+t.numColumns])
		items[rowIndex*t.numColumns+t.numColumns+rowIndex] = column[rowIndex]
	}

	t.items = items
	t.numColumns += 1
}

func (t *Table) At(rowIndex, columnIndex int) *Instance {
	return t.items[rowIndex*t.numColumns+columnIndex]
}

func (t *Table) LabelFromColumnIndex(columnIndex int) *Feature {
	return &t.features[columnIndex]
}

func (t *Table) ColumnIndexFromLabel(typeKey TypeKey) int {
	return t.featureMap[typeKey]
}

func (t *Table) NumColumns() int {
	return t.numColumns
}

func (t *Table) NumRows() int {
	return t.numRows
}
