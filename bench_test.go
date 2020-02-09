package raftsqlite

import (
	"github.com/hashicorp/raft/bench"
	"os"
	"testing"
)

func BenchmarkStore_FirstIndex(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.FirstIndex(b, store)
}

func BenchmarkStore_LastIndex(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.LastIndex(b, store)
}

func BenchmarkStore_GetLog(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.GetLog(b, store)
}

func BenchmarkStore_StoreLog(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.StoreLog(b, store)
}

func BenchmarkStore_StoreLogs(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.StoreLogs(b, store)
}

func BenchmarkStore_DeleteRange(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.DeleteRange(b, store)
}

func BenchmarkStore_Set(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.Set(b, store)
}

func BenchmarkStore_Get(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.Get(b, store)
}

func BenchmarkStore_SetUint64(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.SetUint64(b, store)
}

func BenchmarkStore_GetUint64(b *testing.B) {
	store := testStore(b)
	defer store.Close()
	defer os.Remove(store.path)

	raftbench.GetUint64(b, store)
}
