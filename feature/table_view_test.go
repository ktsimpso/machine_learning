package feature_test

import (
	"github.com/ktsimpso/machine_learning/feature"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

const (
	ITERATIONS = 40
	SUB_TABLES = 15
)

var baseTables []feature.TableViewer

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	baseTables = []feature.TableViewer{
		baseTable(),
	}

	for i := 0; i < SUB_TABLES; i++ {
		baseTables = append(baseTables, baseTableView(baseTables[i]))
	}
}

func baseTable() feature.TableViewer {
	table := feature.CreateTable(featureList)
	table.AddStringRow([]string{"3.14", "2.0", "456.23", "a", "e", "i"})
	table.AddStringRow([]string{"2.44", "1.2", "111.3", "b", "d", "h"})
	table.AddStringRow([]string{"31", "42.0", "44444", "c", "f", "g"})
	table.AddStringRow([]string{"123", "142.1", "55555", "a", "e", "i"})
	table.AddStringRow([]string{"566", "414.0", "66666", "b", "d", "h"})
	table.AddStringRow([]string{"252", "155.5", "77777", "c", "f", "g"})
	table.AddStringRow([]string{"764", "515.3", "88888", "a", "e", "i"})
	table.AddStringRow([]string{"1125", "1516.0", "99999", "b", "d", "h"})
	table.AddStringRow([]string{"1412", "4156.2", "11111", "c", "f", "g"})
	table.AddStringRow([]string{"1241", "1442.8", "22222", "a", "e", "i"})
	table.AddStringRow([]string{"1556", "1245.6", "33333", "b", "d", "h"})
	table.AddStringRow([]string{"63453", "1578.3", "121212", "c", "e", "g"})
	table.AddStringRow([]string{"14124", "12387.4", "232323", "a", "f", "i"})

	return table
}

func baseTableView(table feature.TableViewer) feature.TableViewer {
	builder := feature.NewTableViewBuilder(table)

	for i := 0; i < table.NumRows(); i++ {
		if rand.Float64() < 0.75 {
			builder.WithRow(i)
		}
	}

	for i := 0; i < table.NumColumns(); i++ {
		if rand.Float64() < 0.75 {
			builder.WithColumn(i)
		}
	}

	return builder.Build()
}

func TestViewNumColumns(t *testing.T) {
	assert := assert.New(t)

	for _, table := range baseTables {
		view := feature.NewTableViewBuilder(table).WithAllColumns().WithAllRows().Build()

		assert.Equal(table.NumColumns(), view.NumColumns())

		for iterations := 0; iterations < ITERATIONS; iterations++ {
			count := 0
			builder := feature.NewTableViewBuilder(table).WithAllRows()

			for i := 0; i < table.NumColumns(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithColumn(i)
					count++
				}
			}

			view = builder.Build()
			assert.Equal(count, view.NumColumns())
		}
	}
}

func TestViewLabelFromColumnIndex(t *testing.T) {
	assert := assert.New(t)

	for _, table := range baseTables {
		view := feature.NewTableViewBuilder(table).WithAllColumns().WithAllRows().Build()

		for index := 0; index < table.NumColumns(); index++ {
			assert.Equal(table.LabelFromColumnIndex(index).TypeKey(), view.LabelFromColumnIndex(index).TypeKey())
		}

		for iterations := 0; iterations < ITERATIONS; iterations++ {
			indexes := []int{}
			builder := feature.NewTableViewBuilder(table).WithAllRows()

			for i := 0; i < table.NumColumns(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithColumn(i)
					indexes = append(indexes, i)
				}
			}

			view = builder.Build()
			for index, tableIndex := range indexes {
				assert.Equal(table.LabelFromColumnIndex(tableIndex).TypeKey(), view.LabelFromColumnIndex(index).TypeKey())
			}
		}
	}
}

func TestViewColumnIndexFromLabel(t *testing.T) {
	assert := assert.New(t)

	for _, table := range baseTables {
		view := feature.NewTableViewBuilder(table).WithAllColumns().WithAllRows().Build()

		for index := 0; index < table.NumColumns(); index++ {
			assert.Equal(index, view.ColumnIndexFromLabel(table.LabelFromColumnIndex(index).TypeKey()))
		}

		continue
		for iterations := 0; iterations < ITERATIONS; iterations++ {
			typeKeys := []feature.TypeKey{}
			builder := feature.NewTableViewBuilder(table).WithAllRows()

			for i := 0; i < table.NumColumns(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithColumn(i)
					typeKeys = append(typeKeys, table.LabelFromColumnIndex(i).TypeKey())
				}
			}

			view = builder.Build()
			for index, typeKey := range typeKeys {
				assert.Equal(index, view.ColumnIndexFromLabel(typeKey))
			}
		}
	}
}

func TestViewNumRows(t *testing.T) {
	assert := assert.New(t)

	for _, table := range baseTables {
		view := feature.NewTableViewBuilder(table).WithAllColumns().WithAllRows().Build()

		assert.Equal(table.NumRows(), view.NumRows())

		for iterations := 0; iterations < ITERATIONS; iterations++ {
			count := 0
			builder := feature.NewTableViewBuilder(table).WithAllColumns()

			for i := 0; i < table.NumRows(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithRow(i)
					count++
				}
			}

			view = builder.Build()
			assert.Equal(count, view.NumRows())
		}
	}
}

func TestViewAt(t *testing.T) {
	assert := assert.New(t)

	for _, table := range baseTables {
		view := feature.NewTableViewBuilder(table).WithAllColumns().WithAllRows().Build()

		for rowIndex := 0; rowIndex < table.NumRows(); rowIndex++ {
			for columnIndex := 0; columnIndex < table.NumColumns(); columnIndex++ {
				assert.Equal(table.At(rowIndex, columnIndex).StringValue, view.At(rowIndex, columnIndex).StringValue)
			}
		}

		for iterations := 0; iterations < ITERATIONS; iterations++ {
			rowIndexes := []int{}
			columnIndexes := []int{}
			builder := feature.NewTableViewBuilder(table)

			for i := 0; i < table.NumRows(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithRow(i)
					rowIndexes = append(rowIndexes, i)
				}
			}

			for i := 0; i < table.NumColumns(); i++ {
				if rand.Float64() < 0.5 {
					builder.WithColumn(i)
					columnIndexes = append(columnIndexes, i)
				}
			}

			view = builder.Build()

			for rowIndex, backingRowIndex := range rowIndexes {
				for columnIndex, backingColumnIndex := range columnIndexes {
					assert.Equal(table.At(backingRowIndex, backingColumnIndex), view.At(rowIndex, columnIndex))
				}
			}
		}
	}
}
