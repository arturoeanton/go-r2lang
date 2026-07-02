// Spawns N goroutines that each increment a shared counter under a
// sync.Mutex, using a sync.WaitGroup to know when they have all finished.
let mu = sync.Mutex()
let wg = sync.WaitGroup()

let counter = { n: 0 }
let workers = 20
let incrementsPerWorker = 500

wg.add(workers)

for (let i = 0; i < workers; i = i + 1) {
    r2(func() {
        let j = 0
        while (j < incrementsPerWorker) {
            mu.lock()
            counter.n = counter.n + 1
            mu.unlock()
            j = j + 1
        }
        wg.done()
    })
}

wg.wait()

let expected = workers * incrementsPerWorker
std.print("expected:", expected)
std.print("counter:", counter.n)
if (counter.n == expected) {
    std.print("SYNC_SMOKE_OK")
} else {
    std.print("SYNC_SMOKE_FAIL")
}
