package main

import (
"bufio"
"sims/storage"
)

// SchoolSystem holds shared application services
// used across the portal.
type SchoolSystem struct {
DB             storage.StorageEngine
ConsoleScanner *bufio.Scanner
}
