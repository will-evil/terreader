package terreader

import "fmt"

type dbfSource struct {
	rows []map[string]string
}

func newDbfTable(rows []map[string]string) *dbfSource {
	return &dbfSource{rows: rows}
}

func (ds *dbfSource) NumberOfRecords() int {
	return len(ds.rows)
}

func (ds *dbfSource) FieldValueByName(row int, fieldName string) (string, error) {
	value, ok := ds.rows[row][fieldName]
	if !ok {
		return "", fmt.Errorf("field '%s' not exists", fieldName)
	}

	return value, nil
}
