package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/repo"
	"github.com/fxamacker/cbor/v2"
	cid "github.com/ipfs/go-cid"
)

func main() {
	// Specify the path to your CAR file
	carFilePath := "repo.car"

	// Parse the CAR file and extract records
	err := parseCarFile(carFilePath)
	if err != nil {
		log.Fatalf("Error parsing CAR file: %v", err)
	}
}

// Function to parse the CAR file and extract raw CBOR data
func parseCarFile(carPath string) error {
	ctx := context.Background()

	// Open the CAR file
	fi, err := os.Open(carPath)
	if err != nil {
		return fmt.Errorf("failed to open CAR file: %w", err)
	}
	defer fi.Close()

	// Read the repository from the CAR file
	r, err := repo.ReadRepoFromCar(ctx, fi)
	if err != nil {
		return fmt.Errorf("failed to read repository from CAR: %w", err)
	}

	// Extract DID from the repository commit
	sc := r.SignedCommit()
	did, err := syntax.ParseDID(sc.Did)
	if err != nil {
		return fmt.Errorf("failed to parse DID: %w", err)
	}
	topDir := did.String()

	// Create the output directory for storing raw CBOR files
	err = os.MkdirAll(topDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Access the blockstore from the repo
	blockstore := r.Blockstore()

	// Map to aggregate data by lexicon type
	aggregatedData := make(map[string][]map[string]interface{})

	// Iterate over all records in the repository
	err = r.ForEach(ctx, "", func(k string, v cid.Cid) error {
		// Fetch the raw CBOR block from the blockstore
		block, err := blockstore.Get(ctx, v)
		if err != nil {
			return fmt.Errorf("failed to fetch block for key %s: %w", k, err)
		}

		// Decode the CBOR data
		var decodedData map[interface{}]interface{}
		err = cbor.Unmarshal(block.RawData(), &decodedData)
		if err != nil {
			return fmt.Errorf("failed to decode CBOR for key %s: %w", k, err)
		}

		// Convert map[interface{}]interface{} to map[string]interface{}
		convertedData := convertMapKeysToString(decodedData).(map[string]interface{})

		// Determine the lexicon type from "$type" field
		lexiconType, ok := convertedData["$type"].(string)
		if !ok {
			lexiconType = "unknown" // Fallback if $type is missing
		}

		// Aggregate the record under its lexicon type
		aggregatedData[lexiconType] = append(aggregatedData[lexiconType], convertedData)

		// Save the raw CBOR data
		cborPath := filepath.Join(topDir, k+".cbor")
		err = os.MkdirAll(filepath.Dir(cborPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory for record %s: %w", k, err)
		}
		err = os.WriteFile(cborPath, block.RawData(), 0666)
		if err != nil {
			return fmt.Errorf("failed to write raw block to file %s: %w", cborPath, err)
		}
		fmt.Printf("Raw CBOR record saved: %s\n", cborPath)

		return nil
	})

	if err != nil {
		return fmt.Errorf("error iterating over records: %w", err)
	}

	// Save aggregated data to JSON files
	for lexiconType, records := range aggregatedData {
		fileName := strings.ReplaceAll(lexiconType, ".", "_") + ".json"
		err := saveAggregatedDataAsJSON(records, fileName)
		if err != nil {
			return fmt.Errorf("failed to save aggregated data for %s: %w", lexiconType, err)
		}
	}

	fmt.Printf("All raw records have been saved and aggregated JSON files created.\n")
	return nil
}

// Function to save aggregated data as a JSON file
func saveAggregatedDataAsJSON(data []map[string]interface{}, fileName string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode aggregated JSON: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0666)
	if err != nil {
		return fmt.Errorf("failed to write aggregated JSON file: %w", err)
	}

	fmt.Printf("Aggregated JSON saved: %s\n", fileName)
	return nil
}

// Recursive function to convert map keys to strings
func convertMapKeysToString(input interface{}) interface{} {
	switch v := input.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			newMap[fmt.Sprintf("%v", key)] = convertMapKeysToString(value)
		}
		return newMap
	case []interface{}:
		for i, elem := range v {
			v[i] = convertMapKeysToString(elem)
		}
		return v
	default:
		return input
	}
}
