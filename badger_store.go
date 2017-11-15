package raftbadger

import (
	"os"
	"path"

	"github.com/dgraph-io/badger"
	"github.com/hashicorp/raft"
)

var (
	iterAscOpt  = badger.IteratorOptions{}
	iterDescOpt = badger.IteratorOptions{
		Reverse: true,
	}

	ErrKeyNotFound = badger.ErrKeyNotFound
)

type BadgerOption func(*badger.Options)

func Compact(b bool) BadgerOption {
	return func(o *badger.Options) {
		o.DoNotCompact = !b
	}
}

// TODO: more BadgerOptions

type BadgerStore struct {
	logdb  *badger.DB
	confdb *badger.DB

	dir     string
	logdir  string
	confdir string
}

func open(dir string, opt ...BadgerOption) (*badger.DB, error) {
	o := badger.DefaultOptions
	o.Dir = dir
	o.ValueDir = dir
	o.DoNotCompact = true

	for _, one := range opt {
		one(&o)
	}

	return badger.Open(o)
}

func NewBadgerStore(dir string, opt ...BadgerOption) (*BadgerStore, error) {
	logdir := path.Join(dir, "log")
	confdir := path.Join(dir, "conf")

	if err := os.MkdirAll(logdir, 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(confdir, 0755); err != nil {
		return nil, err
	}

	logdb, err := open(logdir)
	if err != nil {
		return nil, err
	}

	confdb, err := open(confdir)
	if err != nil {
		logdb.Close()
		return nil, err
	}

	return &BadgerStore{
		logdb:   logdb,
		confdb:  confdb,
		dir:     dir,
		logdir:  logdir,
		confdir: confdir,
	}, nil
}

func (b *BadgerStore) Close() error {
	if err := b.logdb.Close(); err != nil {
		return err
	}

	return b.confdb.Close()
}

func (b *BadgerStore) FirstIndex() (uint64, error) {
	tx := b.logdb.NewTransaction(false)
	defer tx.Discard()

	iter := tx.NewIterator(iterAscOpt)
	iter.Rewind()

	item := iter.Item()

	if item == nil {
		return 0, nil
	}

	return bytesToUint64(item.Key()), nil
}

func (b *BadgerStore) LastIndex() (uint64, error) {
	tx := b.logdb.NewTransaction(false)
	defer tx.Discard()

	iter := tx.NewIterator(iterDescOpt)
	iter.Rewind()

	item := iter.Item()

	if item == nil {
		return 0, nil
	}

	return bytesToUint64(item.Key()), nil
}

func (b *BadgerStore) GetLog(idx uint64, log *raft.Log) error {
	tx := b.logdb.NewTransaction(false)
	defer tx.Discard()

	item, err := tx.Get(uint64ToBytes(idx))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return raft.ErrLogNotFound
		}

		return err
	}

	val, err := item.Value()
	if err != nil {
		return err
	}

	return decodeMsgPack(val, log)
}

func (b *BadgerStore) StoreLog(log *raft.Log) error {
	return b.StoreLogs([]*raft.Log{log})
}

func (b *BadgerStore) StoreLogs(logs []*raft.Log) error {
	tx := b.logdb.NewTransaction(true)
	defer tx.Discard()

	for _, one := range logs {
		buf, err := encodeMsgPack(one)
		if err != nil {
			return err
		}

		if err := tx.Set(uint64ToBytes(one.Index), buf.Bytes()); err != nil {
			return err
		}
	}

	if err := tx.Commit(nil); err != nil {
		return err
	}

	return nil
}

func (b *BadgerStore) DeleteRange(min, max uint64) error {
	tx := b.logdb.NewTransaction(true)
	defer tx.Discard()

	minKey := uint64ToBytes(min)

	iter := tx.NewIterator(iterAscOpt)

	for iter.Seek(minKey); iter.Valid(); iter.Next() {
		item := iter.Item()
		if item == nil {
			break
		}

		curKey := safeKey(item)

		if bytesToUint64(curKey) > max {
			break
		}

		if err := tx.Delete(curKey); err != nil {
			return err
		}
	}

	if err := tx.Commit(nil); err != nil {
		return err
	}

	return nil
}

func (b *BadgerStore) Set(key, val []byte) error {
	tx := b.confdb.NewTransaction(true)
	defer tx.Discard()

	if err := tx.Set(key, val); err != nil {
		return err
	}

	if err := tx.Commit(nil); err != nil {
		return err
	}

	return nil
}

func (b *BadgerStore) Get(key []byte) ([]byte, error) {
	tx := b.confdb.NewTransaction(false)
	defer tx.Discard()

	item, err := tx.Get(key)
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

func (b *BadgerStore) SetUint64(key []byte, val uint64) error {
	return b.Set(key, uint64ToBytes(val))
}

func (b *BadgerStore) GetUint64(key []byte) (uint64, error) {
	val, err := b.Get(key)
	if err != nil {
		return 0, err
	}

	return bytesToUint64(val), nil
}

func safeKey(item *badger.Item) []byte {
	key := item.Key()
	dst := make([]byte, len(key))
	copy(dst, key)
	return dst
}
