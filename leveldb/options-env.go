package leveldb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/base/goleveldb/leveldb/opt"
)

func readOptionsFromEnv(options *opt.Options) {
	// ######################################
	// # Reasonable things to try and tune
	// ######################################

	// default is 4 KiB
	setIntFromEnv("LDB_BLOCK_SIZE", &options.BlockSize)

	// ** Compaction Table Size **

	// default is 2 MiB
	tableSizeSet := setIntFromEnv("LDB_COMPACTION_TABLE_SIZE", &options.CompactionTableSize)
	// default is 1
	setFloatFromEnv("LDB_COMPACTION_TABLE_SIZE_MULTIPLIER", &options.CompactionTableSizeMultiplier)

	// ** Compaction Total Size **

	// default is 10 MiB
	totalSizeSet := setIntFromEnv("LDB_COMPACTION_TOTAL_SIZE", &options.CompactionTotalSize)
	// default is 10
	setFloatFromEnv("LDB_COMPACTION_TOTAL_SIZE_MULTIPLIER", &options.CompactionTotalSizeMultiplier)

	// If the table size has been changed but the total size has not explicitly
	// been set, set the total size to 5 times the table size to match the
	// default ratio.
	if tableSizeSet && !totalSizeSet {
		options.CompactionTotalSize = options.CompactionTableSize * 5
	}

	// ######################################
	// # Proceed with caution. Don't set any of these environment variables
	// # unless you know exactly what the implications are for the application
	// # leveraging the LevelDB instance.
	// #
	// # Some of these values will likely be set at runtime by the application
	// # that instantiates the LevelDB instance.
	// ######################################

	// Default is false
	var disableCompression bool
	setBoolFromEnv("LDB_DISABLE_COMPRESSION", &disableCompression)

	if disableCompression {
		options.Compression = opt.NoCompression
	}

	// Default is false
	setBoolFromEnv("LDB_NO_SYNC", &options.NoSync)
	// Default is false
	setBoolFromEnv("LDB_BLOCK_CACHE_EVICT_REMOVED", &options.BlockCacheEvictRemoved)
	// Default is false
	setBoolFromEnv("LDB_DISABLE_BUFFER_POOL", &options.DisableBufferPool)
	// Default is false
	setBoolFromEnv("LDB_DISABLE_BLOCK_CACHE", &options.DisableBlockCache)
	// Default is false
	setBoolFromEnv("LDB_DISABLE_COMPACTION_BACKOFF", &options.DisableCompactionBackoff)
	// Default is false
	setBoolFromEnv("LDB_DISABLE_LARGE_BATCH_TRANSACTION", &options.DisableLargeBatchTransaction)
	// Default is false
	setBoolFromEnv("LDB_DISABLE_SEEKS_COMPACTION", &options.DisableSeeksCompaction)
	// Default is false
	setBoolFromEnv("LDB_ERROR_IF_EXIST", &options.ErrorIfExist)
	// Default is false
	setBoolFromEnv("LDB_ERROR_IF_MISSING", &options.ErrorIfMissing)
	// Default is false
	setBoolFromEnv("LDB_NO_WRITE_MERGE", &options.NoWriteMerge)
	// Default is false
	setBoolFromEnv("LDB_READ_ONLY", &options.ReadOnly)

	// Default is 8 MiB
	setIntFromEnv("LDB_BLOCK_CACHE_CAPACITY", &options.BlockCacheCapacity)
	// Default is 16
	setIntFromEnv("LDB_BLOCK_RESTART_INTERVAL", &options.BlockRestartInterval)
	// Default is 25
	setIntFromEnv("LDB_COMPACTION_EXPAND_LIMIT_FACTOR", &options.CompactionExpandLimitFactor)
	// Default is 10
	setIntFromEnv("LDB_COMPACTION_GP_OVERLAPS_FACTOR", &options.CompactionGPOverlapsFactor)
	// Default is 4
	setIntFromEnv("LDB_COMPACTION_L0_TRIGGER", &options.CompactionL0Trigger)
	// Default is 1
	setIntFromEnv("LDB_COMPACTION_SOURCE_LIMIT_FACTOR", &options.CompactionSourceLimitFactor)
	// Default is 1 MiB
	setIntFromEnv("LDB_ITERATOR_SAMPLING_RATE", &options.IteratorSamplingRate)
	// Default is 4 MiB
	setIntFromEnv("LDB_WRITE_BUFFER", &options.WriteBuffer)
	// Default is 12
	setIntFromEnv("LDB_WRITE_L0_PAUSE_TRIGGER", &options.WriteL0PauseTrigger)
	// Default is 8
	setIntFromEnv("LDB_WRITE_L0_SLOWDOWN_TRIGGER", &options.WriteL0SlowdownTrigger)
	// Default is 11
	setIntFromEnv("LDB_FILTER_BASE_LG", &options.FilterBaseLg)

	// Default is 200 on MacOS and 500 on all other operating systems
	setIntFromEnv("LDB_OPEN_FILES_CACHE_CAPACITY", &options.OpenFilesCacheCapacity)

	var debugOptions bool
	setBoolFromEnv("LDB_DEBUG_OPTIONS", &debugOptions)

	if debugOptions {
		out, _ := json.MarshalIndent(options, "", "    ")
		fmt.Printf("LevelDB Options:\n%s\n", string(out))
	}
}

// `setIntFromEnv` sets the underlying value referenced by pointer `dest` to the parsed value
// of the environment variable with the provided `key`.
//
// If the environment variable is not set or `dest` is nil, the value referenced by `dest` is left
// unchanged and a warning is logged.
func setIntFromEnv(key string, dest *int) (changed bool) {
	strVal := os.Getenv(key)

	if strVal == "" {
		return false
	}

	if dest == nil {
		log.Printf("WARN: Provided destination pointer for parsing of LevelDB env var %s is nil", key)
		log.Print("WARN: Leaving existing value unchanged")
		return false
	}

	initialVal := *dest

	val, err := strconv.Atoi(strVal)
	if err != nil {
		log.Printf("WARN: Error parsing provided LevelDB env var %s to int: %v", key, err)
		log.Printf("WARN: Leaving existing value %d unchanged", *dest)
		return false
	}

	*dest = val
	return initialVal != val
}

func setFloatFromEnv(key string, dest *float64) (changed bool) {
	strVal := os.Getenv(key)

	if strVal == "" {
		return false
	}

	if dest == nil {
		log.Printf("WARN: Provided destination pointer for parsing of LevelDB env var %s is nil", key)
		log.Print("WARN: Leaving existing value unchanged")
		return false
	}

	initialVal := *dest

	val, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Printf("WARN: Error parsing provided LevelDB env var %s to float64: %v", key, err)
		log.Printf("WARN: Leaving existing value %f unchanged", *dest)
		return false
	}

	*dest = val
	return initialVal != val
}

// `setBoolFromEnv` sets the underlying value referenced by pointer `dest` to the parsed value
// of the environment variable with the provided `key`.
//
// If the environment variable is not set or `dest` is nil, the value referenced by `dest` is left
// unchanged and a warning is logged.
func setBoolFromEnv(key string, dest *bool) (changed bool) {
	strVal := os.Getenv(key)

	if strVal == "" {
		return false
	}

	if dest == nil {
		log.Printf("WARN: Provided destination pointer for parsing of LevelDB env var %s is nil", key)
		log.Print("WARN: Leaving existing value unchanged")
		return false
	}

	initialVal := *dest

	val, err := strconv.ParseBool(strVal)
	if err != nil {
		log.Printf("WARN: Error parsing provided LevelDB env var %s to bool: %v", key, err)
		log.Printf("WARN: Leaving existing value %t unchanged", *dest)
		return
	}

	*dest = val
	return initialVal != val
}
