// Package terreader provides functional for reading data from terrorists database which comes as a dbf file.
package terreader

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/will-evil/go-dbf/godbf"
)

type dbfTable interface {
	NumberOfRecords() int
	FieldValueByName(row int, fieldName string) (string, error)
}

// rowData structure for store info about row from dbf terrorist file.
// This struct stores row index in file and value of column ROW-ID.
type rowData struct {
	index int
	rowID uint64
}

type rowDataMap map[uint64][]rowData

// TerReader structure that provides functionality for reading dbf file.
type TerReader struct {
	dbfTable   dbfTable
	rowDataMap rowDataMap
	rowNumbers []uint64
}

// NewTerReader is a constructor for TerReader structure.
func NewTerReader(filePath, encoding string) (*TerReader, error) {
	dbfTable, err := godbf.NewFromFile(filePath, encoding)
	if err != nil {
		return nil, err
	}

	return &TerReader{dbfTable: dbfTable}, nil
}

// Read return chan for retry records from dbf terrorist file.
func (tr *TerReader) Read() (chan RowReadResult, error) {
	if err := tr.setHelpData(); err != nil {
		return nil, err
	}

	rowChan := make(chan RowReadResult, 5)
	go func() {
		resWithError := func(number uint64, err error) RowReadResult {
			return RowReadResult{Number: number, Error: err}
		}

		for _, number := range tr.rowNumbers {
			rowDataSlice, ok := tr.rowDataMap[number]
			if !ok {
				rowChan <- resWithError(number, fmt.Errorf("key '%d' not exists is map rowDataMap", number))
				break
			}

			row, err := tr.buildRecord(rowDataSlice)
			if err != nil {
				rowChan <- resWithError(number, err)
				break
			}

			rowChan <- RowReadResult{Row: row}
		}

		close(rowChan)
	}()

	return rowChan, nil
}

func (tr *TerReader) setHelpData() error {
	if len(tr.rowDataMap) >= 1 {
		return nil
	}

	tr.rowDataMap = make(rowDataMap)

	for i := 0; i < tr.dbfTable.NumberOfRecords(); i++ {
		numberStr, err := tr.dbfTable.FieldValueByName(i, "NUMBER")
		if err != nil {
			return err
		}
		rowIDStr, err := tr.dbfTable.FieldValueByName(i, "ROW_ID")
		if err != nil {
			return err
		}

		number, err := strconv.ParseUint(numberStr, 10, 64)
		if err != nil {
			return err
		}
		rowID, err := strconv.ParseUint(rowIDStr, 10, 64)
		if err != nil {
			return err
		}

		data := rowData{index: i, rowID: rowID}
		if _, ok := tr.rowDataMap[number]; !ok {
			tr.rowNumbers = append(tr.rowNumbers, number)
		}
		tr.rowDataMap[number] = append(tr.rowDataMap[number], data)
	}

	sort.Slice(tr.rowNumbers, func(i, j int) bool {
		return tr.rowNumbers[i] < tr.rowNumbers[j]
	})

	return nil
}

func (tr *TerReader) buildRecord(rowDataSlice []rowData) (*Row, error) {
	if len(rowDataSlice) == 0 {
		return nil, errors.New("rowDataSlice can not be empty")
	}

	sort.SliceStable(rowDataSlice, func(i, j int) bool {
		return rowDataSlice[i].rowID < rowDataSlice[j].rowID
	})

	row := &Row{}

	if err := tr.setStaticFields(row, rowDataSlice[0].index); err != nil {
		return nil, err
	}

	if err := tr.setDateFields(row, rowDataSlice[0].index); err != nil {
		return nil, err
	}

	for _, rowData := range rowDataSlice {
		if err := tr.setEnumFields(row, rowData.index); err != nil {
			return nil, err
		}
		if err := tr.setTextFields(row, rowData.index); err != nil {
			return nil, err
		}
	}

	return row, nil
}

func (tr *TerReader) setStaticFields(row *Row, rowIndex int) error {
	staticFields := []struct {
		name  string
		value *string
	}{
		{"NUMBER", &row.Number},
		{"KODCR", &row.Kodcr},
		{"KODCN", &row.Kodcn},
		{"SD", &row.Sd},
		{"RG", &row.Rg},
		{"ND", &row.Nd},
		{"VD", &row.Vd},
		{"YR", &row.Yr},
		{"ROW_ID", &row.RowID},
	}

	for _, field := range staticFields {
		v, err := tr.dbfTable.FieldValueByName(rowIndex, field.name)
		if err != nil {
			return err
		}
		*field.value = v
	}

	return nil
}

func (tr *TerReader) setDateFields(row *Row, rowIndex int) error {
	getDate := func(fieldName string) (*time.Time, error) {
		str, err := tr.dbfTable.FieldValueByName(rowIndex, fieldName)
		if err != nil {
			return nil, err
		}

		if str == "" {
			return nil, nil
		}

		t, err := time.Parse("20060102", str)

		return &t, err
	}

	grDate, err := getDate("GR")
	if err != nil {
		return err
	}
	row.Gr = grDate

	cbDate, err := getDate("CB_DATE")
	if err != nil {
		return err
	}
	row.CbDate = cbDate

	ceDate, err := getDate("CE_DATE")
	if err != nil {
		return err
	}
	row.CeDate = ceDate

	return nil
}

func (tr *TerReader) setEnumFields(row *Row, rowIndex int) error {
	enumFields := []struct {
		name  string
		enum  []string
		value *string
	}{
		{"TERROR", []string{"0", "1"}, &row.Terror},
		{"TU", []string{"1", "2", "3"}, &row.Tu},
		{"KD", []string{"0", "01", "02", "03", "04"}, &row.Kd},
	}

	isInclude := func(el string, slice []string) bool {
		for _, v := range slice {
			if el == v {
				return true
			}
		}

		return false
	}

	for _, field := range enumFields {
		if isInclude(*field.value, field.enum) {
			continue
		}

		str, err := tr.dbfTable.FieldValueByName(rowIndex, field.name)
		if err != nil {
			return err
		}

		*field.value = str
	}

	return nil
}

func (tr *TerReader) setTextFields(row *Row, rowIndex int) error {
	return nil
}
