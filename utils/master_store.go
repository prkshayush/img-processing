package utils

import (
	"encoding/csv"
	"os"
)

var StoreMaster map[string]bool

func LoadMasterStore(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	StoreMaster = make(map[string]bool)
	for _, record := range records {
		storeID := record[2]
		StoreMaster[storeID] = true
	}

	return nil
}
