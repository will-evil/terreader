package terreader

import (
	"context"
	"errors"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/will-evil/go-dbf/godbf"
)

const fileEncoding = "866"
const filePath = "./test/data/testfile.dbf"

type readTestCase struct {
	rows    []map[string]string
	results []RowReadResult
	err     error
}

func Test_NewTerReader(t *testing.T) {
	testCases := []struct {
		path     string
		encoding string
		reader   *TerReader
		err      error
	}{
		{"", fileEncoding, nil, errors.New("open : no such file or directory")},
		{"not/exists/file.dbf", fileEncoding, nil, errors.New("open not/exists/file.dbf: no such file or directory")},
		{"./test/data/testfile.dbf", fileEncoding, &TerReader{}, nil},
	}

	for _, testCase := range testCases {
		reader, err := NewTerReader(testCase.path, testCase.encoding)
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

		if testCase.reader == nil && reader != nil {
			t.Errorf("get not correct Row. Expecter nil, got %+v", *reader)
		}

		if testCase.reader != nil && reader == nil {
			t.Error("get not correct Row. Expecter structure, got nil")
		}
	}
}

func Test_NewTerReaderFromByteSlice(t *testing.T) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	reader, err := NewTerReaderFromByteSlice(b, fileEncoding)
	if err != nil {
		t.Fatal(err)
	}
	if reader == nil {
		t.Error("get not correct Row. Expecter structure, got nil")
	}
}

func Test_NewTerReaderFromByteSlice_WhenReturnError(t *testing.T) {
	etalonError := errors.New("NewTerReaderFromByteSlice_Error")

	newFromByteSlice = func(data []byte, fileEncoding string) (*godbf.DbfTable, error) {
		return nil, etalonError
	}

	reader, err := NewTerReaderFromByteSlice([]byte{}, fileEncoding)
	if err != etalonError {
		t.Fatalf("error object not correct. Expected %v, got %v", etalonError, err)
	}
	if reader != nil {
		t.Errorf("get not correct Row. Expecter nil, got %v", reader)
	}
}

func Test_WithContext(t *testing.T) {
	tr, err := NewTerReader(filePath, fileEncoding)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tr.WithContext(ctx)

	rowReadRes, err := tr.Read(5)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := (<-rowReadRes); ok {
		t.Error("channel is still open")
	}
}

func Test_TerReader_Read(t *testing.T) {
	for _, testCase := range getTestCases() {
		dbfTable := newDbfTable(testCase.rows)
		tr := TerReader{dbfTable: dbfTable, ctx: context.Background()}

		rowReadRes, err := tr.Read(5)
		if testCase.err == nil && err != nil {
			t.Fatal(err)
		}
		if testCase.err != nil {
			if err == nil {
				t.Fatalf("error object not correct. Expected %v, got nil", testCase.err)
			} else if err.Error() != testCase.err.Error() {
				t.Fatalf("error message not correct. Expected \"%s\", got \"%s\"", testCase.err.Error(), err.Error())
			}
			continue
		}

		num := 0
		etalonRecords := testCase.results
		for res := range rowReadRes {
			etalon := etalonRecords[num]
			if etalon.Error == nil && res.Error != nil {
				t.Fatal(res.Error)
			}
			if etalon.Error != nil {
				if res.Error == nil {
					t.Errorf("error object not correct. Expected %v, got nil", etalon.Error)
				} else if res.Error.Error() != etalon.Error.Error() {
					t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", etalon.Error.Error(), res.Error.Error())
				}
			}

			if res.Number != etalon.Number {
				t.Errorf("Number not correct/ Expected %d, go %d", res.Number, etalon.Number)
			}

			if etalon.Row != nil {
				if res.Row == nil {
					t.Errorf("get not correct Row. Expecter %+v, got nil", *etalon.Row)
				} else if !reflect.DeepEqual(*res.Row, *etalon.Row) {
					t.Errorf("get not correct Row. Expecter %+v, got %+v", *etalon.Row, *res.Row)
				}
			} else if res.Row != nil {
				t.Errorf("get not correct Row. Expecter nil, got %+v", *res.Row)
			}

			num++
		}

		etalonNum := len(etalonRecords)
		if num != etalonNum {
			t.Errorf("received not correct num of records. Expected %d, got %d", etalonNum, num)
		}
	}
}

func Test_TerReader_Reader_WithRealFile(t *testing.T) {
	tr, err := NewTerReader(filePath, fileEncoding)
	if err != nil {
		t.Fatal(err)
	}

	rowReadRes, err := tr.Read(5)
	if err != nil {
		t.Fatal(err)
	}

	Gr := time.Date(1988, time.September, 05, 0, 0, 0, 0, time.UTC)
	CbDate := time.Date(2012, time.November, 10, 0, 0, 0, 0, time.UTC)
	etalonRecords := []RowReadResult{
		{&Row{Number: "1", Terror: "1", Tu: "3", Nameu: "Pharetra magna ac placerat", Descript: "Facilisi etiam dignissim diam quis enim lobortis. Suscipit adipiscing bibendum est ultricies integer quis auctor. At tempor commodo ullamcorper a lacus vestibulum sed arcu. Augue ut lectus arcu bibendum. Porttitor rhoncus dolor purus non enim praesent. Ac tincidunt vitae semper quis lectus nulla at volutpat diam.", Kodcr: "", Kodcn: "", Amr: "", Address: "", Kd: "03", Sd: "", Rg: "BN 5236025", Nd: "", Vd: "Et tortor consequat id porta.", Gr: &Gr, Yr: "1996", Mr: "", CbDate: &CbDate, CeDate: nil, Director: "", Founder: "", RowID: "1", Terrtype: ""}, 1, nil},
	}

	num := 0
	for res := range rowReadRes {
		etalon := etalonRecords[num]
		if etalon.Error == nil && res.Error != nil {
			t.Fatal(res.Error)
		}
		if etalon.Error != nil {
			if res.Error == nil {
				t.Errorf("error object not correct. Expected %v, got nil", etalon.Error)
			} else if res.Error.Error() != etalon.Error.Error() {
				t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", etalon.Error.Error(), res.Error.Error())
			}
		}

		if res.Number != etalon.Number {
			t.Errorf("Number not correct/ Expected %d, go %d", res.Number, etalon.Number)
		}

		if !reflect.DeepEqual(*res.Row, *etalon.Row) {
			t.Errorf("get not correct Row. Expecter %+v, got %+v", *etalon.Row, *res.Row)
		}

		num++
	}

	etalonNum := len(etalonRecords)
	if num != etalonNum {
		t.Errorf("received not correct num uf records. Expected %d, got %d", etalonNum, num)
	}
}

func Test_TerReader_Reader_WhenHelpDataIncorrect(t *testing.T) {
	tr := TerReader{
		rowDataMap: rowDataMap{
			2: []rowData{},
		},
		rowNumbers: []uint64{1},
		ctx:        context.Background(),
	}

	rowReadRes, err := tr.Read(5)
	if err != nil {
		t.Fatal(err)
	}

	etalonError := errors.New("key '1' not exists is map rowDataMap")
	for res := range rowReadRes {
		if res.Row != nil {
			t.Errorf("get not correct Row. Expecter nil, got %+v", *res.Row)
		}
		if res.Error.Error() != etalonError.Error() {
			t.Errorf("error message not correct. Expected \"%s\", got \"%s\"", etalonError.Error(), res.Error.Error())
		}
	}
}

func getTestCases() []readTestCase {
	return []readTestCase{
		{
			rows:    getSuccessTestRows(),
			results: getSuccessEtalonRecords(),
		},
		// When rows is empty.
		{
			rows:    []map[string]string{},
			results: []RowReadResult{},
		},
		// Can not find a suitable value error.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "not_support", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "04", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("can not find a suitable value for 'TERROR'")},
			},
		},
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "not_support", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "04", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("can not find a suitable value for 'TU'")},
			},
		},
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "not_support", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("can not find a suitable value for 'KD'")},
			},
		},
		// Not exists error for static field.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "03", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("field 'KODCR' not exists")},
			},
		},
		// Not exists error for text field.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "04", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("field 'NAMEU' not exists")},
			},
		},
		// Not exists error for enum field.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("field 'KD' not exists")},
			},
		},
		// Not exists error for date field.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "01", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			results: []RowReadResult{
				{Number: 1, Error: errors.New("field 'CB_DATE' not exists")},
			},
		},
		// Not exists error for field 'NUMBER' which needed for help data.
		{
			rows: []map[string]string{
				{"TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "0", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			err: errors.New("field 'NUMBER' not exists"),
		},
		// Not exists error for field 'ROW_ID' which needed for help data.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "0", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "TERRTYPE": ""},
			},
			err: errors.New("field 'ROW_ID' not exists"),
		},
		// Parsing uint error when 'NUMBER' field has not uint value.
		{
			rows: []map[string]string{
				{"NUMBER": "not_int", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "0", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": ""},
			},
			err: errors.New("strconv.ParseUint: parsing \"not_int\": invalid syntax"),
		},
		// Parsing uint error when 'ROW_ID' field has not uint value.
		{
			rows: []map[string]string{
				{"NUMBER": "1", "TERROR": "1", "TU": "1", "NAMEU": "", "DESCRIPT": "", "KODCR": "", "KODCN": "", "AMR": "", "ADRESS": "", "KD": "0", "SD": "", "RG": "", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "-1", "TERRTYPE": ""},
			},
			err: errors.New("strconv.ParseUint: parsing \"-1\": invalid syntax"),
		},
	}
}

func getSuccessTestRows() []map[string]string {
	return []map[string]string{
		{"NUMBER": "1", "TERROR": "1", "TU": "3", "NAMEU": "Olubunmi Pam", "DESCRIPT": "", "KODCR": "004-97", "KODCN": "8624", "AMR": "4258 Queens Lane", "ADRESS": "2066 Confederate Drive", "KD": "01", "SD": "06458975", "RG": "632514", "ND": "718580865862", "VD": "knjoocgbbh", "GR": "20200821", "YR": "2008", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "1", "TERRTYPE": "Resolution 1989"},
		{"NUMBER": "1", "TERROR": "1", "TU": "3", "NAMEU": "Olubunmi Pam", "DESCRIPT": "Diam quam nulla porttitor massa", "KODCR": "004-97", "KODCN": "8624", "AMR": "4258 Queens Lane", "ADRESS": "2066 Confederate Drive", "KD": "01", "SD": "06458975", "RG": "632514", "ND": "718580865862", "VD": "knjoocgbbh", "GR": "20200821", "YR": "2008", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "2", "TERRTYPE": "Resolution 1989"},
		{"NUMBER": "3", "TERROR": "0", "TU": "2", "NAMEU": "Luctus accumsan", "DESCRIPT": "", "KODCR": "004-55", "KODCN": "9632", "AMR": "8855 venenatis", "ADRESS": "624 Venenatis", "KD": "04", "SD": "654684", "RG": "233044", "ND": "46761616", "VD": "pellentesque", "GR": "", "YR": "2005", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "5", "TERRTYPE": "Resolution 1959"},
		{"NUMBER": "2", "TERROR": "1", "TU": "3", "NAMEU": "Neque sodales", "DESCRIPT": "Quis risus sed vulputate odio", "KODCR": "005-66", "KODCN": "3256", "AMR": "4258 Queens Lane", "ADRESS": "2066 Confederate Drive", "KD": "01", "SD": "06458975", "RG": "632514", "ND": "718580865862", "VD": "Interdum posuere", "GR": "20200821", "YR": "", "MR": "Elit pellentesque", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "Adipiscing enim", "ROW_ID": "4", "TERRTYPE": "Sed risus pretium"},
		{"NUMBER": "2", "TERROR": "", "TU": "3", "NAMEU": "Neque sodales", "DESCRIPT": "Quis risus sed vulputate odio", "KODCR": "005-66", "KODCN": "3256", "AMR": "4258 Queens Lane", "ADRESS": "2066 Confederate Drive", "KD": "01", "SD": "06458975", "RG": "632514", "ND": "718580865862", "VD": "Interdum posuere", "GR": "20200821", "YR": "2005", "MR": "Elit pellentesque", "CB_DATE": "20190501", "CE_DATE": "20110101", "DIRECTOR": "Ac felis donec et odio", "FOUNDER": "", "ROW_ID": "3", "TERRTYPE": ""},
		{"NUMBER": "4", "TERROR": "1", "TU": "1", "NAMEU": "Nulla facilisi nullam vehicula ipsum a arcu cursus. Elit eget gravida cum sociis natoque penatibus et magnis. Maecenas volutpat blandit aliquam etiam erat eta. Venenatis lectus magna fringilla urna porttitor rhoncus dolor purus. Fermentum posuere urna ne", "DESCRIPT": "", "KODCR": "654-133", "KODCN": "5238", "AMR": "Neque volutpat", "ADRESS": "", "KD": "01", "SD": "58692315", "RG": "865125", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "5", "TERRTYPE": ""},
		{"NUMBER": "4", "TERROR": "1", "TU": "1", "NAMEU": "c tincidunt praesent semper feugiat nibh.", "DESCRIPT": "", "KODCR": "654-133", "KODCN": "5238", "AMR": "Neque volutpat", "ADRESS": "", "KD": "01", "SD": "58692315", "RG": "865125", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "6", "TERRTYPE": ""},
		{"NUMBER": "5", "TERROR": "1", "TU": "1", "NAMEU": "In cursus turpis massa tincidunt dui ut ornare lectus sit. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Sagittis id consectetur purus ut. Vel pharetra vel turpisu nunc eget lorem dolor sed. Urna id volutpatar lacusan laoreet.", "DESCRIPT": "", "KODCR": "052-752", "KODCN": "2580", "AMR": "Neque volutpat", "ADRESS": "", "KD": "01", "SD": "3688", "RG": "548877", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "7", "TERRTYPE": ""},
		{"NUMBER": "5", "TERROR": "1", "TU": "1", "NAMEU": "Amet facilisis magna etiam tempor orci eu lobortis elementum nibh.", "DESCRIPT": "", "KODCR": "052-752", "KODCN": "2580", "AMR": "Neque volutpat", "ADRESS": "", "KD": "01", "SD": "3688", "RG": "548877", "ND": "", "VD": "", "GR": "", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "8", "TERRTYPE": ""},
		{"NUMBER": "6", "TERROR": "1", "TU": "1", "NAMEU": "Dolor sed viverra ipsum nunc", "DESCRIPT": "Dolor purus non enim praesent. Et pharetra pharetra massa massa ultricies. Fermentum odio eu feugiat pretium. A diam maecenas sed enim ut sem viverra. Duis ut diam quam nulla porttitor massa id neque. Ac tortor dignissim convallis aenean et tortor at ris", "KODCR": "688-888", "KODCN": "6822", "AMR": "", "ADRESS": "", "KD": "04", "SD": "", "RG": "", "ND": "464655", "VD": "", "GR": "20060223", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "9", "TERRTYPE": ""},
		{"NUMBER": "6", "TERROR": "1", "TU": "1", "NAMEU": "Dolor sed viverra ipsum nunc", "DESCRIPT": "us viverra. Viverra nibh cras pulvinar mattis nunc sed blandit.", "KODCR": "688-888", "KODCN": "6822", "AMR": "", "ADRESS": "", "KD": "04", "SD": "", "RG": "", "ND": "464655", "VD": "", "GR": "20060223", "YR": "", "MR": "", "CB_DATE": "", "CE_DATE": "", "DIRECTOR": "", "FOUNDER": "", "ROW_ID": "10", "TERRTYPE": ""},
	}
}

func getSuccessEtalonRecords() []RowReadResult {
	Gr := time.Date(2020, time.August, 21, 0, 0, 0, 0, time.UTC)

	CbDate2 := time.Date(2019, time.May, 01, 0, 0, 0, 0, time.UTC)
	CeDate2 := time.Date(2011, time.January, 01, 0, 0, 0, 0, time.UTC)

	Gr9 := time.Date(2006, time.February, 23, 0, 0, 0, 0, time.UTC)

	return []RowReadResult{
		{&Row{Number: "1", Terror: "1", Tu: "3", Nameu: "Olubunmi Pam", Descript: "Diam quam nulla porttitor massa", Kodcr: "004-97", Kodcn: "8624", Amr: "4258 Queens Lane", Address: "2066 Confederate Drive", Kd: "01", Sd: "06458975", Rg: "632514", Nd: "718580865862", Vd: "knjoocgbbh", Gr: &Gr, Yr: "2008", Mr: "", CbDate: nil, CeDate: nil, Director: "", Founder: "", RowID: "1", Terrtype: "Resolution 1989"}, 1, nil},
		{&Row{Number: "2", Terror: "1", Tu: "3", Nameu: "Neque sodales", Descript: "Quis risus sed vulputate odio", Kodcr: "005-66", Kodcn: "3256", Amr: "4258 Queens Lane", Address: "2066 Confederate Drive", Kd: "01", Sd: "06458975", Rg: "632514", Nd: "718580865862", Vd: "Interdum posuere", Gr: &Gr, Yr: "2005", Mr: "Elit pellentesque", CbDate: &CbDate2, CeDate: &CeDate2, Director: "Ac felis donec et odio", Founder: "Adipiscing enim", RowID: "3", Terrtype: "Sed risus pretium"}, 2, nil},
		{&Row{Number: "3", Terror: "0", Tu: "2", Nameu: "Luctus accumsan", Descript: "", Kodcr: "004-55", Kodcn: "9632", Amr: "8855 venenatis", Address: "624 Venenatis", Kd: "04", Sd: "654684", Rg: "233044", Nd: "46761616", Vd: "pellentesque", Gr: nil, Yr: "2005", Mr: "", CbDate: nil, CeDate: nil, Director: "", Founder: "", RowID: "5", Terrtype: "Resolution 1959"}, 3, nil},
		{&Row{Number: "4", Terror: "1", Tu: "1", Nameu: "Nulla facilisi nullam vehicula ipsum a arcu cursus. Elit eget gravida cum sociis natoque penatibus et magnis. Maecenas volutpat blandit aliquam etiam erat eta. Venenatis lectus magna fringilla urna porttitor rhoncus dolor purus. Fermentum posuere urna nec tincidunt praesent semper feugiat nibh.", Descript: "", Kodcr: "654-133", Kodcn: "5238", Amr: "Neque volutpat", Address: "", Kd: "01", Sd: "58692315", Rg: "865125", Nd: "", Vd: "", Gr: nil, Yr: "", Mr: "", CbDate: nil, CeDate: nil, Director: "", Founder: "", RowID: "5", Terrtype: ""}, 4, nil},
		{&Row{Number: "5", Terror: "1", Tu: "1", Nameu: "In cursus turpis massa tincidunt dui ut ornare lectus sit. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Sagittis id consectetur purus ut. Vel pharetra vel turpisu nunc eget lorem dolor sed. Urna id volutpatar lacusan laoreet. Amet facilisis magna etiam tempor orci eu lobortis elementum nibh.", Descript: "", Kodcr: "052-752", Kodcn: "2580", Amr: "Neque volutpat", Address: "", Kd: "01", Sd: "3688", Rg: "548877", Nd: "", Vd: "", Gr: nil, Yr: "", Mr: "", CbDate: nil, CeDate: nil, Director: "", Founder: "", RowID: "7", Terrtype: ""}, 5, nil},
		{&Row{Number: "6", Terror: "1", Tu: "1", Nameu: "Dolor sed viverra ipsum nunc", Descript: "Dolor purus non enim praesent. Et pharetra pharetra massa massa ultricies. Fermentum odio eu feugiat pretium. A diam maecenas sed enim ut sem viverra. Duis ut diam quam nulla porttitor massa id neque. Ac tortor dignissim convallis aenean et tortor at risus viverra. Viverra nibh cras pulvinar mattis nunc sed blandit.", Kodcr: "688-888", Kodcn: "6822", Amr: "", Address: "", Kd: "04", Sd: "", Rg: "", Nd: "464655", Vd: "", Gr: &Gr9, Yr: "", Mr: "", CbDate: nil, CeDate: nil, Director: "", Founder: "", RowID: "9", Terrtype: ""}, 6, nil},
	}
}
