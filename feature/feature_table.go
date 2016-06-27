package feature

type Table struct {
	Features   []Feature
	FeatureMap map[TypeKey]int
	Columns    [][]*Instance
	Rows       [][]*Instance
}

func CreateTable(features []Feature) *Table {
	featureMap := map[TypeKey]int{}

	for index, feature := range features {
		featureMap[feature.TypeKey()] = index
	}

	return &Table{
		features,
		featureMap,
		make([][]*Instance, len(features)),
		[][]*Instance{},
	}
}

func (t *Table) AddStringRow(records []string) {
	//TOOD: error checking
	row := make([]*Instance, len(t.Features))
	for index, feature := range t.Features {
		instance := feature.Create(records[index])

		row[index] = instance
		t.Columns[index] = append(t.Columns[index], instance)
	}

	t.Rows = append(t.Rows, row)
}

func (t *Table) AddRow(row []*Instance) {
	//TODO: error checking
	for index := range t.Features {
		t.Columns[index] = append(t.Columns[index], row[index])
	}

	t.Rows = append(t.Rows, row)
}

func (t *Table) AddColumn(feature Feature, column []*Instance) {
	//TODO: error checking
	t.Features = append(t.Features, feature)
	t.Columns = append(t.Columns, column)
	t.FeatureMap[feature.TypeKey()] = len(t.Features) - 1

	for index, row := range t.Rows {
		t.Rows[index] = append(row, column[index])
	}
}
