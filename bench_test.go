package raftbadger

import (
	"os"
	"testing"

	"github.com/hashicorp/raft/bench"
)

func BenchmarkBadgerStore_FirstIndex(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.FirstIndex(b, store)
}

func BenchmarkBadgerStore_LastIndex(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.LastIndex(b, store)
}

func BenchmarkBadgerStore_GetLog(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.GetLog(b, store)
}

func BenchmarkBadgerStore_StoreLog(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.StoreLog(b, store)
}

func BenchmarkBadgerStore_StoreLogs(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.StoreLogs(b, store)
}

func BenchmarkBadgerStore_DeleteRange(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.DeleteRange(b, store)
}

func BenchmarkBadgerStore_Set(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.Set(b, store)
}

func BenchmarkBadgerStore_Get(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.Get(b, store)
}

func BenchmarkBadgerStore_SetUint64(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.SetUint64(b, store)
}

func BenchmarkBadgerStore_GetUint64(b *testing.B) {
	store := testBadgerStore(b)
	defer store.Close()
	defer os.Remove(store.dir)

	raftbench.GetUint64(b, store)
}
