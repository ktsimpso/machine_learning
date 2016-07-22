package feature

type Table struct {
	features   []Feature
	featureMap map[TypeKey]int
	rows       [][]*Instance
}

func CreateTable(features []Feature) *Table {
	table := Table{}
	featureMap := map[TypeKey]int{}

	for index, feature := range features {
		featureMap[feature.TypeKey()] = index
	}

	table.features = features
	table.featureMap = featureMap
	table.rows = [][]*Instance{}

	return &table
}

func (t *Table) AddStringRow(records []string) {
	//TOOD: error checking
	row := make([]*Instance, len(t.features))

	for index, feature := range t.features {
		row[index] = feature.Create(records[index])
	}

	t.rows = append(t.rows, row)
}

func (t *Table) AddColumn(feature Feature, column []*Instance) {
	//TODO: error checking
	t.features = append(t.features, feature)
	t.featureMap[feature.TypeKey()] = len(t.features) - 1

	for index, _ := range t.rows {
		t.rows[index] = append(t.rows[index], column[index])
	}
}

func (t *Table) At(rowIndex, columnIndex int) *Instance {
	return t.rows[rowIndex][columnIndex]
}

func (t *Table) LabelFromColumnIndex(columnIndex int) *Feature {
	return &t.features[columnIndex]
}

func (t *Table) ColumnIndexFromLabel(typeKey TypeKey) int {
	return t.featureMap[typeKey]
}

func (t *Table) NumColumns() int {
	return len(t.features)
}

func (t *Table) NumRows() int {
	return len(t.rows)
}
