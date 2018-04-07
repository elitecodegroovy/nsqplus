package threads

import (
	"testing"
	"fmt"
	"time"
	"runtime"
)

func TestSimpleCase(t *testing.T) {

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs+1) // numCPUs hot threads + one for async tasks.
	//--Creating a thread pool, with 30 executors and max of 1 million pending jobs
	var thread_pool *ThreadPool = NewPool(numCPUs * 10, 1000000)
	// Note that if pending job is more than 1000000, the new submission (call to Submit) will be blocked
	// until the job queue has some space.

	// thread_pool.Open() must be called. Without this, threads won't start processing jobs
	thread_pool.Open()
	// After thread_pool is started, there is 30 go routines in background, processing jobs

	//---Submiting a job and gets a Future
	var fut *Future = thread_pool.Submit(func() interface{} {
		var j int = 0
		for i := 0; i < 10000; i++ {
			j += i*i
		}
		return  j
	})

	//Getting stats of this pool
	fmt.Println(" ActiveCount: ", thread_pool.ActiveCount()) // active jobs - being executed right now
	fmt.Println(" PendingCount: ", thread_pool.PendingCount()) // pending count - not started yet
	fmt.Println(" CompletedCount: ", thread_pool.CompletedCount()) //jobs done - result populated already

	// Here, submited func that returns a value. The func will be executed by a backend processor
	// where there is free go routine. The submission returns a *threads.Future, which can be used
	// to retrieve the returned value from the func.
	// e.g.

	//----Wait until the future is ready to be retrieve
	result := fut.GetWait().(int) // <= result will be 
	fmt.Println("Result of 1 + 6 is", result)
	// Wait until it is run and result is ready

	// or if you prefer no blocking, call returns immediately, but may contain no result
	ok, result1 := fut.GetNoWait()
	if ok {
		// result is ready
		fmt.Println("Result of 1 + 6 is", result1) // <= result will be 7
	} else {
		fmt.Println("Result is not ready yet")
	}

	//Getting stats of this pool
	fmt.Println(" ActiveCount: ", thread_pool.ActiveCount()) // active jobs - being executed right now
	fmt.Println(" PendingCount: ", thread_pool.PendingCount()) // pending count - not started yet
	fmt.Println(" CompletedCount: ", thread_pool.CompletedCount()) //jobs done - result populated already

	// or if you want to wait for max 3 seconds
	ok, result2 := fut.GetWaitTimeout(3*time.Second)
	if ok {
	// result is ready
		fmt.Println("Result of 1 + 6 is", result2) // <= result will be 7
	} else {
		fmt.Println("Result is not ready yet") // <= timed out after 3 seconds
	}

	//-- Stop accepting new jobs
	// once shutdown, you can't re-start it back
	thread_pool.Shutdown()
	// Now thread_pool can't submit new jobs. All existing submited jobs will be still processed
	// The future previous returned will still materialize

	// Wait until all jobs to complete. Calling Wait() on non-shutdown thread pool will be blocked forever
	thread_pool.Wait()
	// after this call, all futures should be able to be retrieved without delay
	// You can safely disregard this thread_pool after this call. It is useless anyway

	//Getting stats of this pool
	thread_pool.ActiveCount() // active jobs - being executed right now
	thread_pool.PendingCount() // pending count - not started yet
	thread_pool.CompletedCount() //jobs done - result populated already

	//--Convenient wrapper to do multiple tasks in parallel
	//jobs := make([]func() interface{}, 60)
	////... populate the jobs with actual jobs
	//// This will start as many threads as possible to run things in parallel
	//var fg FutureGroup = ParallelDo(jobs)
	//
	//// This will start at most 10 threads for parallel processing
	//var fg FutureGroup = ParallelDoWithLimit(jobs, 10)
	//
	//// retrieve futures, wait for all and get result!
	//var results[]interface{} = fg.WaitAll()
	//
	//// If you prefer more flexible handling...
	//var []*threads.Future futures = fg.Futures
}
