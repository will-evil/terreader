package terreader

import "time"

// Row is struct for store data of row from file.
type Row struct {
	Number   uint64
	Terror   bool
	Tu       uint8
	Nameu    string
	Descript string
	Kodcr    string
	Kodcn    string
	Amr      string
	Address  string
	Kd       string
	Sd       string
	Rg       string
	Nd       string
	Vd       string
	Gr       *time.Time
	Yr       string
	Mr       string
	CbDate   *time.Time
	CeDate   *time.Time
	Director string
	Founder  string
	RowID    string
	Terrtype string
}
