package hw06_pipeline_execution //nolint:golint,stylecheck
import (
	"sort"
	"sync"
)

type (
	In          = <-chan interface{}
	Out         = In
	Bi          = chan interface{}
	St          = chan StageResult
	StageResult struct {
		WorkerIndex int
		Value       interface{}
	}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	values := []interface{}{}
	for val := range in {
		values = append(values, val)
	}
	isDone := false
	mutex := sync.Mutex{}

	go func() {
		<-done
		mutex.Lock()
		isDone = true
		mutex.Unlock()
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(values))
	stageResCh := make(St, len(values))

	for i, val := range values {
		go func(workerIndex int, val interface{}) {
			defer wg.Done()
			for _, stage := range stages {
				mutex.Lock()
				if isDone {
					mutex.Unlock()
					return
				}
				mutex.Unlock()

				inCh := make(Bi, 1)
				defer close(inCh)

				inCh <- val
				val = <-stage(inCh)
			}

			stageResCh <- StageResult{WorkerIndex: workerIndex, Value: val}
		}(i, val)
	}
	wg.Wait()
	close(stageResCh)

	stagingResults := []StageResult{}
	for r := range stageResCh {
		stagingResults = append(stagingResults, r)
	}

	sort.Slice(stagingResults, func(i, j int) bool {
		return stagingResults[i].WorkerIndex < stagingResults[j].WorkerIndex
	})

	resCh := make(Bi, len(values))
	defer close(resCh)

	for _, r := range stagingResults {
		resCh <- r.Value
	}

	return resCh
}
