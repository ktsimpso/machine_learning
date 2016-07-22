package feature

type TableViewer interface {
	NumColumns() int
	NumRows() int

	At(rowIndex, columnIndex int) *Instance
	LabelFromColumnIndex(columnIndex int) *Feature
	ColumnIndexFromLabel(typeKey TypeKey) int
}

type TableViewBuilder struct {
	backingView   TableViewer
	columnIndexes []int
	rowIndexes    []int
}

func NewTableViewBuilder(table TableViewer) *TableViewBuilder {
	return &TableViewBuilder{
		table,
		[]int{},
		[]int{},
	}
}

func (tvb *TableViewBuilder) WithAllColumns() *TableViewBuilder {
	tvb.columnIndexes = make([]int, tvb.backingView.NumColumns())

	for i := 0; i < len(tvb.columnIndexes); i++ {
		tvb.columnIndexes[i] = i
	}

	return tvb
}

func (tvb *TableViewBuilder) WithAllRows() *TableViewBuilder {
	tvb.rowIndexes = make([]int, tvb.backingView.NumRows())

	for i := 0; i < len(tvb.rowIndexes); i++ {
		tvb.columnIndexes[i] = i
	}

	return tvb
}

func (tvb *TableViewBuilder) WithColumn(index int) *TableViewBuilder {
	tvb.columnIndexes = append(tvb.columnIndexes, index)
	return tvb
}

func (tvb *TableViewBuilder) WithRow(index int) *TableViewBuilder {
	tvb.rowIndexes = append(tvb.rowIndexes, index)
	return tvb
}

func (tvb *TableViewBuilder) Build() TableViewer {
	featureMap := map[TypeKey]int{}

	switch backing := tvb.backingView.(type) {
	case *Table:
		for index, backingIndex := range tvb.columnIndexes {
			featureMap[backing.features[backingIndex].TypeKey()] = index
		}

		return &TableView{
			backing,
			featureMap,
			tvb.columnIndexes,
			tvb.rowIndexes,
		}
	case *TableView:
		columnIndexes := make([]int, len(tvb.columnIndexes))
		featureMap := map[TypeKey]int{}

		for index, backingIndex := range tvb.columnIndexes {
			columnIndexes[index] = backing.columnIndexes[backingIndex]
			featureMap[backing.backingTable.features[backingIndex].TypeKey()] = index
		}

		rowIndexes := make([]int, len(tvb.rowIndexes))
		for index, backingIndex := range tvb.rowIndexes {
			rowIndexes[index] = backing.rowIndexes[backingIndex]
		}

		return &TableView{
			backing.backingTable,
			featureMap,
			columnIndexes,
			rowIndexes,
		}
	}

	panic("Not a supported backing TableViewer")
}

type TableView struct {
	backingTable  *Table
	featureMap    map[TypeKey]int
	columnIndexes []int
	rowIndexes    []int
}

func (tv *TableView) NumColumns() int {
	return len(tv.columnIndexes)
}

func (tv *TableView) NumRows() int {
	return len(tv.rowIndexes)
}

func (tv *TableView) At(rowIndex, columnIndex int) *Instance {
	return tv.backingTable.rows[tv.rowIndexes[rowIndex]][tv.columnIndexes[columnIndex]]
}

func (tv *TableView) LabelFromColumnIndex(columnIndex int) *Feature {
	return &tv.backingTable.features[tv.columnIndexes[columnIndex]]
}

func (tv *TableView) ColumnIndexFromLabel(typeKey TypeKey) int {
	return tv.featureMap[typeKey]
}
