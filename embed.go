package vnprovince

import "embed"

// DivisionPath is the path to the divisions.csv file.
const DivisionPath = "data/divisions.csv"

// DataDirFS is the filesystem containing the data directory.
//
//go:embed all:data
var DataDirFS embed.FS
