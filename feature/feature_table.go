package feature

type TableColumn struct {
	instances []*Instance
	table     *Table
	index     int
}

type TableRow struct {
	instances []*Instance
	table     *Table
	index     int
}

type Table struct {
	features   []Feature
	featureMap map[TypeKey]int
	columns    []*TableColumn
	rows       []*TableRow
}

func (tc *TableColumn) At(index int) *Instance {
	return tc.instances[index]
}

func (tc *TableColumn) Feature() Feature {
	return tc.table.features[tc.index]
}

func (tc *TableColumn) Instances() <-chan Record {
	ch := make(chan Record)

	go func() {
		defer close(ch)
		for index, instance := range tc.instances {
			ch <- Record{
				instance,
				index,
				tc.table.features[tc.index],
			}
		}
	}()

	return ch
}

func (tc *TableColumn) Len() int {
	return len(tc.instances)
}

func (tr *TableRow) AtType(typeKey TypeKey) *Instance {
	return tr.instances[tr.table.featureMap[typeKey]]
}

func (tr *TableRow) Index() int {
	return tr.index
}

func (tr *TableRow) Instances() <-chan Record {
	ch := make(chan Record)

	go func() {
		defer close(ch)
		for index, instance := range tr.instances {
			ch <- Record{
				instance,
				index,
				tr.table.features[index],
			}
		}
	}()

	return ch
}

func (tr *TableRow) Len() int {
	return len(tr.instances)
}

func CreateTable(features []Feature) *Table {
	table := Table{}
	featureMap := map[TypeKey]int{}
	columns := make([]*TableColumn, len(features))

	for index, feature := range features {
		featureMap[feature.TypeKey()] = index
		columns[index] = &TableColumn{
			[]*Instance{},
			&table,
			index,
		}
	}

	table.features = features
	table.featureMap = featureMap
	table.columns = columns
	table.rows = []*TableRow{}

	return &table
}

func (t *Table) AddStringRow(records []string) {
	//TOOD: error checking
	row := TableRow{
		make([]*Instance, len(t.features)),
		t,
		len(t.rows),
	}
	for index, feature := range t.features {
		instance := feature.Create(records[index])

		row.instances[index] = instance
		t.columns[index].instances = append(t.columns[index].instances, instance)
	}

	t.rows = append(t.rows, &row)
}

func (t *Table) AddRow(row *TableRow) {
	//TODO: error checking
	tableRow := TableRow{
		row.instances,
		t,
		len(t.rows),
	}
	for index := range t.features {
		t.columns[index].instances = append(t.columns[index].instances, row.instances[index])
	}

	t.rows = append(t.rows, &tableRow)
}

func (t *Table) AddColumn(feature Feature, column []*Instance) {
	//TODO: error checking
	t.features = append(t.features, feature)
	t.columns = append(t.columns, &TableColumn{
		column,
		t,
		len(t.features) - 1,
	})
	t.featureMap[feature.TypeKey()] = len(t.features) - 1

	for index, row := range t.rows {
		//TODO: does the feature and featuremap get updated here?
		row.instances = append(row.instances, column[index])
	}
}

func (t *Table) GetColumn(typeKey TypeKey) Column {
	return t.columns[t.featureMap[typeKey]]
}

func (t *Table) Columns() <-chan Column {
	ch := make(chan Column)

	go func() {
		defer close(ch)
		for _, column := range t.columns {
			ch <- column
		}
	}()

	return ch
}

func (t *Table) NumColumns() int {
	return len(t.columns)
}

func (t *Table) GetRow(index int) Row {
	return t.rows[index]
}

func (t *Table) Rows() <-chan Row {
	ch := make(chan Row)

	go func() {
		defer close(ch)
		for _, row := range t.rows {
			ch <- row
		}
	}()

	return ch
}

func (t *Table) NumRows() int {
	return len(t.rows)
}
