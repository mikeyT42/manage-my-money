package models

import (
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

var logE = log.New(os.Stderr, "BankStatement: ", log.Lshortfile|log.LstdFlags)

type BankStatement struct {
	Transactions []BankTransaction
}

// -----------------------------------------------------------------------------
func Unmarshal(data []byte, bs *BankStatement) error {
	err := gocsv.UnmarshalBytes(data, &bs.Transactions)
	if err != nil {
		logE.Println("Could not unmarshal.")
		logE.Println(err)
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
func (bs *BankStatement) AutoCategorize() {
	d_to_c := buildDescriptionCategoryMap()
    const NONE int32 = 1

	for i := range bs.Transactions {
		transaction := &bs.Transactions[i]
        transaction.CategoryId = NONE
		desc_tokens := strings.Split(transaction.Description, " ")

		for _, token := range desc_tokens {
			category, exists := d_to_c[token]
			if exists {
				transaction.CategoryId = category
				break
			}
		}
	}
}

// -----------------------------------------------------------------------------
func buildDescriptionCategoryMap() map[string]int32 {
	return map[string]int32{
		// Default
		"LABONNE'S":        13,
		"SHOPRITE":         13,
		"*HISVINE":         17,
		"ORTHOCONNECTICUT": 29,
		"3GTMS":            29,
		"JPMorgan":         12,
		"CINTI":            18,
		"T-MOBILE":         14,
		"Spectrum":         15,
		"EVERSOURCE":       16,
		"ALLTOWN":          10,
		"BRIDGEWATER":      10,
		"EXXON":            10,
		"WHEELS":           10,
		"GULF":             10,
		"CUMBERLAND":       10,
		"7-ELEVEN":         10,
        "AMG":              10,
		"CHARLEYS":         11,
		"RESTAURANT":       11,
		"WENDY'S":          11,
		"*AYLA'S":          11,
		"TACO":             11,
		"DUNKIN":           11,
		"McDonalds":        11,
		"STARBUCKS":        11,
		"SONIC":            11,
		"MCDONALD'S":       11,
		"BAGELMAN":         11,
		"SHAKE":            11,
		"MIX":              11,
		"DELI":             11,
		// Custom
		"PL*GWManagement": 26,
		"PL*PAYLEASE":     26,
		"BUTTHEADS":       20,
		"BARNES":          19,
		"PLANET":          22,
		"AUDIOCHUCK.CO":   22,
		"NETFLIX":         22,
		"Prime*RA":        22,
		"MOHELA":          23,
		"MYCHART":         21,
		"APRIA-MEDICAL":   27,
	}
}
