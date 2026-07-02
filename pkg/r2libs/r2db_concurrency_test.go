package r2libs

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// getDBFunc is a small helper to pull a builtin out of the "db" module.
func getDBFunc(t *testing.T, env *r2core.Environment, name string) r2core.BuiltinFunction {
	t.Helper()
	dbModuleObj, ok := env.Get("db")
	if !ok {
		t.Fatal("db module not found")
	}
	dbModule := dbModuleObj.(map[string]interface{})
	fnVal, ok := dbModule[name]
	if !ok {
		t.Fatalf("%s function not found", name)
	}
	return fnVal.(r2core.BuiltinFunction)
}

// TestDBConcurrentConnectNoDuplicateIDs hammers dbConnect from many goroutines
// simultaneously and verifies every returned connection id is unique, then
// closes them all concurrently too.
func TestDBConcurrentConnectNoDuplicateIDs(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbClose := getDBFunc(t, env, "dbClose")

	const n = 100
	ids := make([]string, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			id := dbConnect("sqlite3", ":memory:").(string)
			ids[i] = id
		}(i)
	}
	wg.Wait()

	seen := make(map[string]bool, n)
	for i, id := range ids {
		if id == "" {
			t.Fatalf("goroutine %d got empty connection id", i)
		}
		if seen[id] {
			t.Fatalf("duplicate connection id %q returned to two different goroutines", id)
		}
		seen[id] = true
	}

	// Close everything concurrently.
	var wg2 sync.WaitGroup
	wg2.Add(n)
	for _, id := range ids {
		go func(id string) {
			defer wg2.Done()
			dbClose(id)
		}(id)
	}
	wg2.Wait()
}

// TestDBConcurrentOpenCloseChurnNoCollision opens/closes many connections in a
// tight concurrent loop (simulating heavy churn on dbConnCounter) and
// verifies that no two live connections ever share an id, even though ids
// keep advancing past the current map size.
func TestDBConcurrentOpenCloseChurnNoCollision(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbClose := getDBFunc(t, env, "dbClose")

	const workers = 50
	const itersPerWorker = 20

	var mu sync.Mutex
	allIssued := make(map[string]int) // id -> count of times issued

	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < itersPerWorker; i++ {
				id := dbConnect("sqlite3", ":memory:").(string)

				mu.Lock()
				allIssued[id]++
				mu.Unlock()

				// Immediately close roughly half of them to create churn
				// (freeing map slots) while others stay open.
				if i%2 == 0 {
					dbClose(id)
				}
			}
		}()
	}
	wg.Wait()

	for id, count := range allIssued {
		if count != 1 {
			t.Fatalf("connection id %q was issued %d times concurrently (collision)", id, count)
		}
	}
}

// TestDBOutOfOrderCloseNoStaleReuse opens A, B, C; closes B; opens D; and
// verifies D's id never collides with A or C, and that using B's old id
// afterwards fails with a clean "not found" panic from our own bookkeeping
// rather than silently operating on a reused/stale id or blowing up with a
// raw database/sql use-after-close panic.
func TestDBOutOfOrderCloseNoStaleReuse(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbClose := getDBFunc(t, env, "dbClose")
	dbQuery := getDBFunc(t, env, "dbQuery")

	idA := dbConnect("sqlite3", ":memory:").(string)
	idB := dbConnect("sqlite3", ":memory:").(string)
	idC := dbConnect("sqlite3", ":memory:").(string)

	dbClose(idB)

	idD := dbConnect("sqlite3", ":memory:").(string)

	if idD == idA || idD == idC || idD == idB {
		t.Fatalf("new connection id %q collided with an existing id (A=%q B=%q C=%q)", idD, idA, idB, idC)
	}

	// Using B's old id must panic cleanly with "not found" (our map lookup),
	// not a raw sql.ErrConnDone-style crash and not silently succeed against
	// a reused slot.
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic when querying a closed/removed connection id, got none")
			}
			msg := fmt.Sprintf("%v", r)
			if !strings.Contains(msg, "not found") {
				t.Fatalf("expected a clean 'not found' panic for stale id, got: %v", msg)
			}
		}()
		dbQuery(idB, "SELECT 1")
	}()

	// A, C, D should still be fully usable.
	for _, id := range []string{idA, idC, idD} {
		res := dbQuery(id, "SELECT 1 as one")
		rows := res.([]interface{})
		if len(rows) != 1 {
			t.Fatalf("expected 1 row from surviving connection %q, got %d", id, len(rows))
		}
	}

	dbClose(idA)
	dbClose(idC)
	dbClose(idD)
}

// TestDBConcurrentQueryDuringClose exercises the TOCTOU window in dbQuery /
// dbExec: they RLock to fetch the *sql.DB and RUnlock before actually
// running the query, so a concurrent dbClose can close the underlying
// *sql.DB while a query is in flight or about to start. This verifies that
// outcome is always either a clean success or a recoverable Go panic
// (surfacing database/sql's "database is closed" error), and never a data
// race (-race) or an unrecoverable fatal crash of the test binary.
func TestDBConcurrentQueryDuringClose(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbExec := getDBFunc(t, env, "dbExec")
	dbQuery := getDBFunc(t, env, "dbQuery")
	dbClose := getDBFunc(t, env, "dbClose")

	id := dbConnect("sqlite3", ":memory:").(string)
	dbExec(id, "CREATE TABLE t (id INTEGER PRIMARY KEY, v TEXT)")
	dbExec(id, "INSERT INTO t (v) VALUES (?)", "hello")

	const readers = 30
	var wg sync.WaitGroup
	wg.Add(readers + 1)

	panics := make([]interface{}, readers)

	for i := 0; i < readers; i++ {
		go func(i int) {
			defer wg.Done()
			defer func() {
				panics[i] = recover()
			}()
			for j := 0; j < 20; j++ {
				dbQuery(id, "SELECT * FROM t")
			}
		}(i)
	}

	go func() {
		defer wg.Done()
		dbClose(id)
	}()

	wg.Wait()

	// Nothing should have escalated beyond a normal recoverable panic; if we
	// reach this line at all, the process didn't crash. Sanity check that
	// any panic we did see is the expected database-closed style error and
	// not something bizarre (e.g. a nil pointer dereference indicating a
	// real bug).
	for i, p := range panics {
		if p == nil {
			continue
		}
		msg := fmt.Sprintf("%v", p)
		if !strings.Contains(msg, "dbQuery:") {
			t.Fatalf("reader %d got unexpected panic shape: %v", i, msg)
		}
	}
}

// TestDBConcurrentDoubleClose calls dbClose on the same connection id from
// two goroutines simultaneously. database/sql's *sql.DB.Close is documented
// as safe for concurrent use, so this must not race or panic in a way that
// isn't a clean "not found" (the second delete simply finds nothing).
func TestDBConcurrentDoubleClose(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbClose := getDBFunc(t, env, "dbClose")

	id := dbConnect("sqlite3", ":memory:").(string)

	var wg sync.WaitGroup
	results := make([]interface{}, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer wg.Done()
			defer func() {
				results[i] = recover()
			}()
			dbClose(id)
		}(i)
	}
	wg.Wait()

	// At most one of the two should have panicked (the loser finding the id
	// already gone); neither should be some unexpected error type.
	panicCount := 0
	for _, r := range results {
		if r != nil {
			panicCount++
			msg := fmt.Sprintf("%v", r)
			if !strings.Contains(msg, "not found") && !strings.Contains(msg, "dbClose:") {
				t.Fatalf("unexpected double-close panic shape: %v", msg)
			}
		}
	}
	if panicCount > 1 {
		t.Fatalf("expected at most 1 of 2 concurrent dbClose calls to panic, got %d", panicCount)
	}
}

// TestDBBeginIsDecorativeOnly documents (and locks in) the current, known
// limitation: dbBegin begins a transaction and immediately rolls it back,
// returning a placeholder id that is never stored anywhere and can never be
// committed against. This is pre-existing behavior, not something this test
// pass changes.
func TestDBBeginIsDecorativeOnly(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDB(env)

	dbConnect := getDBFunc(t, env, "dbConnect")
	dbExec := getDBFunc(t, env, "dbExec")
	dbBegin := getDBFunc(t, env, "dbBegin")
	dbQuery := getDBFunc(t, env, "dbQuery")
	dbClose := getDBFunc(t, env, "dbClose")

	id := dbConnect("sqlite3", ":memory:").(string)
	dbExec(id, "CREATE TABLE t (id INTEGER PRIMARY KEY, v TEXT)")

	txId := dbBegin(id).(string)
	if !strings.HasPrefix(txId, "tx_") {
		t.Fatalf("expected tx placeholder id to start with tx_, got %v", txId)
	}

	// The txId is not usable for anything: it was never stored, so using it
	// as a connection id must fail with "not found", proving there is no
	// commit/rollback/exec path wired to it.
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected using the tx placeholder id as a connection id to panic")
			}
		}()
		dbQuery(txId, "SELECT 1")
	}()

	dbClose(id)
}
