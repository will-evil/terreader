// Package terreader provides functional for reading data from terrorists database which comes as a dbf file.
package terreader

import (
	"sort"
	"strconv"

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

// Read return chan for retry rows from dbf terrorist file.
func (tr *TerReader) Read() (chan RowReadResult, error) {
	if err := tr.setHelpData(); err != nil {
		return nil, err
	}

	rowChan := make(chan RowReadResult, 5)
	go func() {
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
