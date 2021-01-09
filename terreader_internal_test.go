package terreader

import (
	"errors"
	"reflect"
	"testing"
)

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

// When rowDataMap already set.
func Test_TerReader_setHelpData_WhenDataMapAlreadySet(t *testing.T) {
	dataMap := rowDataMap{
		1: []rowData{},
	}
	tr := TerReader{rowDataMap: dataMap}
	if err := tr.setHelpData(); err != nil {
		t.Fatal(err)
	}
}

func Test_TerReader_buildRecord(t *testing.T) {
	tr := TerReader{}
	row, err := tr.buildRecord([]rowData{})
	if row != nil {
		t.Errorf("row not correct. Expected nil, got %+v", row)
	}

	etalonError := errors.New("rowDataSlice can not be empty")
	if err.Error() != etalonError.Error() {
		t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", etalonError.Error(), err.Error())
	}
}

func Test_TerReader_getEnumValue(t *testing.T) {
	testCases := []struct {
		fieldName    string
		rowDataSlice []rowData
		res          string
		err          error
	}{
		{"NOT_SUPPORT", []rowData{}, "", errors.New("not support field name 'NOT_SUPPORT'")},
	}

	tr := TerReader{}
	for _, testCase := range testCases {
		res, err := tr.getEnumValue(testCase.fieldName, testCase.rowDataSlice)

		if testCase.err == nil && err != nil {
			t.Fatal(err)
		}
		if testCase.err != nil {
			if err == nil {
				t.Errorf("error object not correct. Expected %v, got nil", testCase.err)
			} else if err.Error() != testCase.err.Error() {
				t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", testCase.err.Error(), err.Error())
			}
		}

		if testCase.res != res {
			t.Errorf("get not correct string. Expected \"%s\", got \"%s\"", testCase.res, res)
		}
	}
}

func Test_getEnum(t *testing.T) {
	testCases := [4]struct {
		field string
		enum  []string
		err   error
	}{
		{"TERROR", []string{"0", "1"}, nil},
		{"TU", []string{"1", "2", "3"}, nil},
		{"KD", []string{"0", "01", "02", "03", "04"}, nil},
		{"NOT_SUPPORT", []string{}, errors.New("not support field name 'NOT_SUPPORT'")},
	}

	for _, testCase := range testCases {
		res, err := getEnum(testCase.field)
		if testCase.err == nil && err != nil {
			t.Fatal(err)
		}
		if testCase.err != nil {
			if err == nil {
				t.Errorf("error object not correct. Expected %v, got nil", testCase.err)
			} else if err.Error() != testCase.err.Error() {
				t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", testCase.err.Error(), err.Error())
			}
		}

		if !reflect.DeepEqual(testCase.enum, res) {
			t.Errorf("not correct result. Expected %v, got %v", testCase.enum, res)
		}
	}
}
