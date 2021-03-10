// Copyright © 2021 Alexey Konovalenko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package terreader provides functional for reading data from terrorists database which comes as a dbf file.
package terreader

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/will-evil/go-dbf/godbf"
)

const (
	dateFormat               = "20060102"
	needTrailingSpaceCharNum = 253
)

var newFromByteSlice = godbf.NewFromByteArray

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

// NewTerReaderFromByteSlice is a TerReader constructor for slice of bytes.
func NewTerReaderFromByteSlice(data []byte, encoding string) (*TerReader, error) {
	dbfTable, err := newFromByteSlice(data, encoding)
	if err != nil {
		return nil, err
	}

	return &TerReader{dbfTable: dbfTable}, nil
}

// Read return chan for retry records from dbf terrorist file.
func (tr *TerReader) Read(chanBuff uint) (chan RowReadResult, error) {
	if err := tr.setHelpData(); err != nil {
		return nil, err
	}

	rowChan := make(chan RowReadResult, chanBuff)
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

			rowChan <- RowReadResult{Row: row, Number: number}
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

	val := reflect.ValueOf(row).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		fieldName := typeField.Tag.Get("tr_col")
		fieldType := typeField.Tag.Get("tr_type")
		switch fieldType {
		case "static":
			val, err := tr.dbfTable.FieldValueByName(rowDataSlice[0].index, fieldName)
			if err != nil {
				return nil, err
			}
			valueField.SetString(val)
		case "enum":
			val, err := tr.getEnumValue(fieldName, rowDataSlice)
			if err != nil {
				return nil, err
			}
			valueField.SetString(val)
		case "date":
			val, err := tr.getDateValue(fieldName, rowDataSlice[0].index)
			if err != nil {
				return nil, err
			}
			valueField.Set(reflect.ValueOf(val))
		case "text":
			val, err := tr.getTextValue(fieldName, rowDataSlice)
			if err != nil {
				return nil, err
			}
			valueField.SetString(val)
		}
	}

	return row, nil
}

func (tr *TerReader) getDateValue(fieldName string, rowIndex int) (*time.Time, error) {
	val, err := tr.dbfTable.FieldValueByName(rowIndex, fieldName)
	if err != nil {
		return nil, err
	}

	if val == "" {
		return nil, nil
	}

	t, err := time.Parse(dateFormat, val)

	return &t, err
}

func (tr *TerReader) getEnumValue(fieldName string, rowDataSlice []rowData) (string, error) {
	enumValues, err := getEnum(fieldName)
	if err != nil {
		return "", err
	}

	isInclude := func(el string, slice []string) bool {
		for _, v := range slice {
			if el == v {
				return true
			}
		}

		return false
	}

	for _, data := range rowDataSlice {
		val, err := tr.dbfTable.FieldValueByName(data.index, fieldName)
		if err != nil {
			return "", err
		}

		if isInclude(val, enumValues) {
			return val, nil
		}
	}

	return "", fmt.Errorf("can not find a suitable value for '%s'", fieldName)
}

func (tr *TerReader) getTextValue(fieldName string, rowDataSlice []rowData) (string, error) {
	var text string
	var lastIncludedStrLen int

	for _, data := range rowDataSlice {
		val, err := tr.dbfTable.FieldValueByName(data.index, fieldName)
		if err != nil {
			return "", err
		}

		if strings.Contains(text, val) || val == "" {
			continue
		}

		leadingChar := ""
		if lastIncludedStrLen == needTrailingSpaceCharNum {
			leadingChar = " "
		}

		text += leadingChar + val

		lastIncludedStrLen = len(val)
	}

	return text, nil
}

func getEnum(fieldName string) ([]string, error) {
	switch fieldName {
	case "TERROR":
		return []string{"0", "1"}, nil
	case "TU":
		return []string{"1", "2", "3"}, nil
	case "KD":
		return []string{"0", "01", "02", "03", "04"}, nil
	}

	return []string{}, fmt.Errorf("not support field name '%s'", fieldName)
}
