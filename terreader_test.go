package terreader

import (
	"reflect"
	"testing"
)

func Test_TerReader_Read(t *testing.T) {

}

func Test_TerReader_setHelpData(t *testing.T) {
	rows := []map[string]string{
		{"NUMBER": "1", "ROW_ID": "2"},
		{"NUMBER": "2", "ROW_ID": "5"},
		{"NUMBER": "2", "ROW_ID": "7"},
		{"NUMBER": "2", "ROW_ID": "6"},
		{"NUMBER": "1", "ROW_ID": "1"},
		{"NUMBER": "1", "ROW_ID": "4"},
		{"NUMBER": "1", "ROW_ID": "3"},
	}
	dbfTable := newDbfTable(rows)
	tr := TerReader{dbfTable: dbfTable}
	tr.setHelpData()

	rowDataMap := map[uint64][]rowData{
		1: {
			{index: 0, rowID: 2},
			{index: 4, rowID: 1},
			{index: 5, rowID: 4},
			{index: 6, rowID: 3},
		},
		2: {
			{index: 1, rowID: 5},
			{index: 2, rowID: 7},
			{index: 3, rowID: 6},
		},
	}
	for number, etalonSlice := range rowDataMap {
		resultSlice, ok := tr.rowDataMap[number]
		if !ok {
			t.Errorf("key '%d' not exists in result map", number)
			continue
		}
		if !reflect.DeepEqual(resultSlice, etalonSlice) {
			t.Errorf("value for key '%d' not correct. Expected %v, got %v", number, etalonSlice, resultSlice)
		}
	}

	rowNumbers := []uint64{1, 2}
	if !reflect.DeepEqual(tr.rowNumbers, rowNumbers) {
		t.Errorf("rowNumbers not equal. Expected %v, got %v", rowNumbers, tr.rowNumbers)
	}
}
