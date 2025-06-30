package vnprovince

import _ "embed"


// DataDirFS is the filesystem containing the data directory.
//
//go:embed data/divisions_16_10_2024.csv
var dataDirFS string
