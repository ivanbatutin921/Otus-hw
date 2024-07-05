package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	// 1)defer close(out)

	go func() {
		OutChan := in

		for _, stage := range stages {
			OutChan = stage(OutChan)
		}

		for v := range OutChan {
			select {
			case out <- v:
			case <-done:
				// 2)close(out)
				return
			}
		}
		// 2)close(out)
	}()
	return out
}
