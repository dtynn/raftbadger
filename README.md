### raftbadger
Raft backend implementation using [Badger](https://github.com/dgraph-io/badger)

#### benchmarks
On my MacbookAir (11-inch, Early 2014)

```
BenchmarkBadgerStore_FirstIndex-4         300000          4494 ns/op
BenchmarkBadgerStore_LastIndex-4          300000          4743 ns/op
BenchmarkBadgerStore_GetLog-4             300000          3814 ns/op
BenchmarkBadgerStore_StoreLog-4            10000        180024 ns/op
BenchmarkBadgerStore_StoreLogs-4            5000        343946 ns/op
BenchmarkBadgerStore_DeleteRange-4          5000        250139 ns/op
BenchmarkBadgerStore_Set-4                 10000        175970 ns/op
BenchmarkBadgerStore_Get-4               1000000          1146 ns/op
BenchmarkBadgerStore_SetUint64-4           10000        182825 ns/op
BenchmarkBadgerStore_GetUint64-4         1000000          1273 ns/op
```


```
BenchmarkBoltStore_FirstIndex-4      1000000          1081 ns/op
BenchmarkBoltStore_LastIndex-4       1000000          1016 ns/op
BenchmarkBoltStore_GetLog-4           500000          3494 ns/op
BenchmarkBoltStore_StoreLog-4           3000        340762 ns/op
BenchmarkBoltStore_StoreLogs-4          3000        421707 ns/op
BenchmarkBoltStore_DeleteRange-4        5000        341560 ns/op
BenchmarkBoltStore_Set-4                5000        321252 ns/op
BenchmarkBoltStore_Get-4             1000000          1209 ns/op
BenchmarkBoltStore_SetUint64-4          5000        312685 ns/op
BenchmarkBoltStore_GetUint64-4       1000000          1205 ns/op
```
