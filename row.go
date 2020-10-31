package terreader

import "time"

// Row is struct for store data of row from file.
type Row struct {
	Number   string     `tr_col:"NUMBER"   tr_type:"static"`
	Terror   string     `tr_col:"TERROR"   tr_type:"enum"`
	Tu       string     `tr_col:"TU"       tr_type:"enum"`
	Nameu    string     `tr_col:"NAMEU"    tr_type:"text"`
	Descript string     `tr_col:"DESCRIPT" tr_type:"text"`
	Kodcr    string     `tr_col:"KODCR"    tr_type:"static"`
	Kodcn    string     `tr_col:"KODCN"    tr_type:"static"`
	Amr      string     `tr_col:"AMR"      tr_type:"text"`
	Address  string     `tr_col:"ADRESS"   tr_type:"text"`
	Kd       string     `tr_col:"KD"       tr_type:"enum"`
	Sd       string     `tr_col:"SD"       tr_type:"static"`
	Rg       string     `tr_col:"RG"       tr_type:"static"`
	Nd       string     `tr_col:"ND"       tr_type:"static"`
	Vd       string     `tr_col:"VD"       tr_type:"static"`
	Gr       *time.Time `tr_col:"GR"       tr_type:"date"`
	Yr       string     `tr_col:"YR"       tr_type:"static"`
	Mr       string     `tr_col:"MR"       tr_type:"text"`
	CbDate   *time.Time `tr_col:"CB_DATE"  tr_type:"date"`
	CeDate   *time.Time `tr_col:"CE_DATE"  tr_type:"date"`
	Director string     `tr_col:"DIRECTOR" tr_type:"text"`
	Founder  string     `tr_col:"FOUNDER"  tr_type:"text"`
	RowID    string     `tr_col:"ROW_ID"   tr_type:"static"`
	Terrtype string     `tr_col:"TERRTYPE" tr_type:"text"`
}
