package vnprovince

import "embed"

const DivisionPath = "data/divisions.csv"

// DataDirFS is the filesystem containing the data directory.
//
//go:embed all:data
var DataDirFS embed.FS
