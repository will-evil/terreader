package terreader

// RowReadResult structure for store result of row reading.
type RowReadResult struct {
	Row    Row
	Number uint64
	Error  error
}
