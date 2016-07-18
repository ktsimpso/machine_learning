package feature

type TableViewer interface {
	GetColumn(typeKey TypeKey) Column
	GetRow(index int) Row

	Columns() <-chan Column
	Rows() <-chan Row

	NumColumns() int
	NumRows() int
}

type Record struct {
	Instance *Instance
	Index    int
	Feature  Feature
}

type Column interface {
	At(index int) *Instance
	Feature() Feature
	Instances() <-chan Record
	Len() int
}

type Row interface {
	AtType(typeKey TypeKey) *Instance
	Index() int
	Instances() <-chan Record
	Len() int
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
	switch backing := tvb.backingView.(type) {
	case *Table:
		return &TableView{
			backing,
			tvb.columnIndexes,
			tvb.rowIndexes,
		}
	case *TableView:
		columnIndexes := make([]int, len(tvb.columnIndexes))
		for index, backingIndex := range tvb.columnIndexes {
			columnIndexes[index] = backing.columnIndexes[backingIndex]
		}

		rowIndexes := make([]int, len(tvb.rowIndexes))
		for index, backingIndex := range tvb.rowIndexes {
			rowIndexes[index] = backing.rowIndexes[backingIndex]
		}

		return &TableView{
			backing.backingTable,
			columnIndexes,
			rowIndexes,
		}
	}

	panic("Not a supported backing TableViewer")
}

type TableView struct {
	backingTable  *Table
	columnIndexes []int
	rowIndexes    []int
}

func (tv *TableView) NumColumns() int {
	return len(tv.columnIndexes)
}

func (tv *TableView) NumRows() int {
	return len(tv.rowIndexes)
}

func (tv *TableView) GetColumn(typeKey TypeKey) Column {
	//TODO: make this efficient
	targetTypeIndex := tv.backingTable.featureMap[typeKey]
	for index, targetIndex := range tv.columnIndexes {
		if targetTypeIndex == targetIndex {
			return &TableViewColumn{
				tv,
				index,
			}
		}
	}

	panic("No column in view with type key")
}

func (tv *TableView) GetRow(index int) Row {
	return &TableViewRow{
		tv,
		index,
	}
}

func (tv *TableView) Columns() <-chan Column {
	ch := make(chan Column)

	go func() {
		defer close(ch)
		for index := range tv.columnIndexes {
			ch <- &TableViewColumn{
				tv,
				index,
			}
		}
	}()

	return ch
}

func (tv *TableView) Rows() <-chan Row {
	ch := make(chan Row)

	go func() {
		defer close(ch)
		for index := range tv.rowIndexes {
			ch <- &TableViewRow{
				tv,
				index,
			}
		}
	}()

	return ch
}

type TableViewColumn struct {
	tableView *TableView
	index     int
}

func (tvc *TableViewColumn) At(index int) *Instance {
	return tvc.tableView.backingTable.columns[tvc.tableView.columnIndexes[tvc.index]].instances[tvc.tableView.rowIndexes[index]]
}

func (tvc *TableViewColumn) Feature() Feature {
	return tvc.tableView.backingTable.features[tvc.tableView.columnIndexes[tvc.index]]
}

func (tvc *TableViewColumn) Instances() <-chan Record {
	ch := make(chan Record)

	go func() {
		defer close(ch)
		backingColumn := tvc.tableView.backingTable.columns[tvc.tableView.columnIndexes[tvc.index]]
		for index, backingIndex := range tvc.tableView.rowIndexes {
			ch <- Record{
				backingColumn.instances[backingIndex],
				index,
				tvc.tableView.backingTable.features[backingColumn.index],
			}
		}
	}()

	return ch
}

func (tvc *TableViewColumn) Len() int {
	return tvc.tableView.NumRows()
}

type TableViewRow struct {
	tableView *TableView
	index     int
}

func (tvr *TableViewRow) AtType(typeKey TypeKey) *Instance {
	//TOOD: error checking
	return tvr.tableView.backingTable.rows[tvr.tableView.rowIndexes[tvr.index]].instances[tvr.tableView.backingTable.featureMap[typeKey]]
}

func (tvr *TableViewRow) Index() int {
	return tvr.index
}

func (tvr *TableViewRow) Instances() <-chan Record {
	ch := make(chan Record)

	go func() {
		defer close(ch)
		backingRow := tvr.tableView.backingTable.rows[tvr.tableView.rowIndexes[tvr.index]]
		for index, backingIndex := range tvr.tableView.columnIndexes {
			ch <- Record{
				backingRow.instances[backingIndex],
				index,
				tvr.tableView.backingTable.features[backingIndex],
			}
		}
	}()

	return ch
}

func (tvr *TableViewRow) Len() int {
	return tvr.tableView.NumColumns()
}
