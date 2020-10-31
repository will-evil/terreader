package terreader

import (
	"reflect"
	"testing"
	"time"
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
	if err := tr.setHelpData(); err != nil {
		t.Fatal(err)
	}

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

func Test_TerReader_setStaticFields(t *testing.T) {
	rows := []map[string]string{
		{
			"NUMBER": "1",
			"KODCR":  "788",
			"KODCN":  "434",
			"SD":     "2601",
			"RG":     "BN",
			"ND":     "060400763505",
			"VD":     "PASSPORT",
			"YR":     "1986",
			"ROW_ID": "2",
		},
	}
	dbfTable := newDbfTable(rows)
	tr := TerReader{dbfTable: dbfTable}

	row := &Row{}
	if err := tr.setStaticFields(row, 0); err != nil {
		t.Fatal(err)
	}

	fieldMap := make(map[string]string)
	fieldMap["NUMBER"] = row.Number
	fieldMap["KODCR"] = row.Kodcr
	fieldMap["KODCN"] = row.Kodcn
	fieldMap["SD"] = row.Sd

	for name, value := range fieldMap {
		if value != rows[0][name] {
			t.Errorf("field '%s' not correct. Expected %s, got %s", name, rows[0][name], value)
		}
	}
}

func Test_TerReader_setDateFields(t *testing.T) {
	rows := []map[string]string{
		{"GR": "20020517", "CB_DATE": "20200126", "CE_DATE": ""},
	}
	dbfTable := newDbfTable(rows)
	tr := TerReader{dbfTable: dbfTable}

	row := &Row{}
	if err := tr.setDateFields(row, 0); err != nil {
		t.Fatal(err)
	}

	fieldMap := make(map[string]*time.Time)
	fieldMap["GR"] = row.Gr
	fieldMap["CB_DATE"] = row.CbDate
	fieldMap["CE_DATE"] = row.CeDate

	var etalonDate *time.Time
	for name, value := range fieldMap {
		if timeStr := rows[0][name]; timeStr != "" {
			time, err := time.Parse("20060102", timeStr)
			if err != nil {
				t.Fatal(err)
			}
			etalonDate = &time
		} else {
			etalonDate = nil
		}

		if value == nil {
			if value != etalonDate {
				t.Errorf("field '%s' not correct. Expected %s, got %s", name, etalonDate, value)
			}
			continue
		}

		if !etalonDate.Equal(*value) {
			t.Errorf("field '%s' not correct. Expected %s, got %s", name, etalonDate, value)
		}
	}
}

func Test_TerReader_setEnumFields(t *testing.T) {
	rows := []map[string]string{
		{"TERROR": "", "TU": "", "KD": "01"},
		{"TERROR": "1", "TU": "2", "KD": "02"},
	}
	dbfTable := newDbfTable(rows)
	tr := TerReader{dbfTable: dbfTable}

	row := &Row{}

	testCases := [2][]struct {
		field  string
		val    *string
		etalon string
	}{
		{
			{"TERROR", &row.Terror, ""},
			{"TU", &row.Tu, ""},
			{"KD", &row.Kd, "01"},
		},
		{
			{"TERROR", &row.Terror, "1"},
			{"TU", &row.Tu, "2"},
			{"KD", &row.Kd, "01"},
		},
	}

	for index, testCase := range testCases {
		if err := tr.setEnumFields(row, index); err != nil {
			t.Fatal(err)
		}

		for _, test := range testCase {
			if *test.val != test.etalon {
				t.Errorf("field '%s' not correct. Expected %s, got %s", test.field, test.etalon, *test.val)
			}
		}
	}
}
