Fork Notice
-----------
Forked from [github.com/syndtr/goleveldb](https://github.com/syndtr/goleveldb) on 2025-03-10 at commit 2ae1ddf74ef7020251ff1ff0fe8daac21a157761.

This corresponds to the version specified as a dependency by go-ethereum and op-geth when this fork was created (github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7).


This fork adds support for setting database initialization parameters via environment variables.

Supported environment variables:

| Name                                          | Type    | Default Value   | Corresponds To                              |
| --------------------------------------------- | ------- | --------------- | ------------------------------------------- |
| `LDB_BLOCK_SIZE`                              | Integer | `4096`          | `opt.Options.BlockSize`                         |
| `LDB_COMPACTION_TABLE_SIZE`                   | Integer | `2097152`       | `opt.Options.CompactionTableSize`               |
| `LDB_COMPACTION_TOTAL_SIZE`                   | Integer | `10485760`      | `opt.Options.CompactionTotalSize`               |
| `LDB_COMPACTION_TOTAL_SIZE_MULTIPLIER`        | Integer | `10`            | `opt.Options.CompactionTotalSizeMultiplier`     |
| `LDB_BLOCK_CACHE_CAPACITY`                    | Integer | `8388608`       | `opt.Options.BlockCacheCapacity`                |
| `LDB_BLOCK_RESTART_INTERVAL`                  | Integer | `16`            | `opt.Options.BlockRestartInterval`              |
| `LDB_COMPACTION_EXPAND_LIMIT_FACTOR`          | Integer | `25`            | `opt.Options.CompactionExpandLimitFactor`       |
| `LDB_COMPACTION_GP_OVERLAPS_FACTOR`           | Integer | `10`            | `opt.Options.CompactionGPOverlapsFactor`        |
| `LDB_COMPACTION_L0_TRIGGER`                   | Integer | `4`             | `opt.Options.CompactionL0Trigger`               |
| `LDB_COMPACTION_SOURCE_LIMIT_FACTOR`          | Integer | `1`             | `opt.Options.CompactionSourceLimitFactor`       |
| `LDB_ITERATOR_SAMPLING_RATE`                  | Integer | `1048576`       | `opt.Options.IteratorSamplingRate`              |
| `LDB_WRITE_BUFFER`                            | Integer | `4194304`       | `opt.Options.WriteBuffer`                       |
| `LDB_WRITE_L0_PAUSE_TRIGGER`                  | Integer | `12`            | `opt.Options.WriteL0PauseTrigger`               |
| `LDB_WRITE_L0_SLOWDOWN_TRIGGER`               | Integer | `8`             | `opt.Options.WriteL0SlowdownTrigger`            |
| `LDB_FILTER_BASE_LG`                          | Integer | `11`            | `opt.Options.FilterBaseLg`                      |
| `LDB_OPEN_FILES_CACHE_CAPACITY`               | Integer | `200` (MacOS), `500` (others) | `opt.Options.OpenFilesCacheCapacity` |
| `LDB_DISABLE_COMPRESSION`                     | Boolean | `false`         | `opt.Options.Compression` (if true, set to `opt.NoCompression`) |
| `LDB_NO_SYNC`                                 | Boolean | `false`         | `opt.Options.NoSync`                            |
| `LDB_BLOCK_CACHE_EVICT_REMOVED`               | Boolean | `false`         | `opt.Options.BlockCacheEvictRemoved`            |
| `LDB_DISABLE_BUFFER_POOL`                     | Boolean | `false`         | `opt.Options.DisableBufferPool`                 |
| `LDB_DISABLE_BLOCK_CACHE`                     | Boolean | `false`         | `opt.Options.DisableBlockCache`                 |
| `LDB_DISABLE_COMPACTION_BACKOFF`              | Boolean | `false`         | `opt.Options.DisableCompactionBackoff`          |
| `LDB_DISABLE_LARGE_BATCH_TRANSACTION`         | Boolean | `false`         | `opt.Options.DisableLargeBatchTransaction`      |
| `LDB_DISABLE_SEEKS_COMPACTION`                | Boolean | `false`         | `opt.Options.DisableSeeksCompaction`            |
| `LDB_ERROR_IF_EXIST`                          | Boolean | `false`         | `opt.Options.ErrorIfExist`                      |
| `LDB_ERROR_IF_MISSING`                        | Boolean | `false`         | `opt.Options.ErrorIfMissing`                    |
| `LDB_NO_WRITE_MERGE`                          | Boolean | `false`         | `opt.Options.NoWriteMerge`                      |
| `LDB_READ_ONLY`                               | Boolean | `false`         | `opt.Options.ReadOnly`                          |
| `LDB_DEBUG_OPTIONS`                           | Boolean | `false`         | will print `opt.Options` to the console after initialization |


-----------

This is an implementation of the [LevelDB key/value database](http:code.google.com/p/leveldb) in the [Go programming language](http:golang.org).


Installation
-----------

	go get github.com/base/goleveldb/leveldb

Requirements
-----------

* Need at least `go1.5` or newer.

Usage
-----------

Create or open a database:
```go
// The returned DB instance is safe for concurrent use. Which mean that all
// DB's methods may be called concurrently from multiple goroutine.
db, err := leveldb.OpenFile("path/to/db", nil)
...
defer db.Close()
...
```
Read or modify the database content:
```go
// Remember that the contents of the returned slice should not be modified.
data, err := db.Get([]byte("key"), nil)
...
err = db.Put([]byte("key"), []byte("value"), nil)
...
err = db.Delete([]byte("key"), nil)
...
```

Iterate over database content:
```go
iter := db.NewIterator(nil, nil)
for iter.Next() {
	// Remember that the contents of the returned slice should not be modified, and
	// only valid until the next call to Next.
	key := iter.Key()
	value := iter.Value()
	...
}
iter.Release()
err = iter.Error()
...
```
Seek-then-Iterate:
```go
iter := db.NewIterator(nil, nil)
for ok := iter.Seek(key); ok; ok = iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```
Iterate over subset of database content:
```go
iter := db.NewIterator(&util.Range{Start: []byte("foo"), Limit: []byte("xoo")}, nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```
Iterate over subset of database content with a particular prefix:
```go
iter := db.NewIterator(util.BytesPrefix([]byte("foo-")), nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```
Batch writes:
```go
batch := new(leveldb.Batch)
batch.Put([]byte("foo"), []byte("value"))
batch.Put([]byte("bar"), []byte("another value"))
batch.Delete([]byte("baz"))
err = db.Write(batch, nil)
...
```
Use bloom filter:
```go
o := &opt.Options{
	Filter: filter.NewBloomFilter(10),
}
db, err := leveldb.OpenFile("path/to/db", o)
...
defer db.Close()
...
```
