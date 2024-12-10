package vnprovince

import "embed"

// DivisionPath is the path to the divisions.csv file.
const DivisionPath = "data/divisions_16_10_2024.csv"

// DataDirFS is the filesystem containing the data directory.
//
//go:embed data/divisions_16_10_2024.csv
var DataDirFS embed.FS
