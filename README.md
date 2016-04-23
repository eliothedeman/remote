# remote
Remote access to boltdb files via rpc

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
So far only basic operations are covered.

### DB
- [ ] Batch
- [x] Begin
- [x] Close
- [ ] GoString
- [ ] Info
- [x] IsReadOnly
- [x] Path
- [ ] Stats
- [ ] String
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
- [ ] ForEach
- [x] Get
- [ ] NextSequence
- [x] Put
- [ ] Root
- [ ] Stats
- [ ] Tx
- [ ] Writeable

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
- [ ] OnCommit
- [ ] Page
- [x] Rollback
- [ ] Size
- [ ] Writeable
- [ ] WriteTo
