# Remote [![GoDoc](https://godoc.org/github.com/eliothedeman/remote?status.svg)](https://godoc.org/github.com/eliothedeman/remote)
Remote access to [boltdb](https://github.com/boltdb/bolt) files via rpc

## Mission
Expose the same great api of boltdb over the network.

## Use

Just like boltdb... but with a network address.
```go
db, err := remote.Open("tcp://10.0.0.1:9090")
```

Or open a local database.
```go
db, err := remote.Open("/home/user/local.db")
```

## Progress
Most common opperations are covered. Certain functions that are done through io
interfaces have been put off until a better streaming system can be implemented.

### DB
- [ ] Batch
- [x] Begin
- [x] Close
- [x] GoString
- [ ] Info (Not applicable)
- [x] IsReadOnly
- [x] Path
- [x] Stats
- [x] String
- [ ] Sync
- [x] Update
- [x] View

### Bucket
- [x] Bucket
- [x] CreateBucket
- [x] CreateBucketIfNotExists
- [ ] Cursor
- [x] Delete
- [x] DeleteBucket
- [x] ForEach
- [x] Get
- [x] NextSequence
- [x] Put
- [ ] Root
- [x] Stats
- [x] Tx
- [x] Writeable

### Tx
- [x] Bucket
- [x] Commit
- [x] CreateBucket
- [x] CreateBucketIfNotExists
- [ ] Copy
- [ ] CopyFile
- [ ] Cursor
- [x] DB
- [x] DeleteBucket
- [ ] ForEach
- [x] OnCommit
- [ ] Page (Not applicable)
- [x] Rollback
- [x] Size
- [ ] Writeable
- [ ] WriteTo

## Plan
Get the api down using the built in go rpc system, then make the move to gRPC to support other language client libraries.
