package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	inCopy := make(Bi)

	tmp := in
	go func() {
		for _, stage := range stages {
			tmp = stage(inCopy)
		}

		for v := range tmp {
			inCopy <- v
		}

	}()

	go func() {
		defer close(out)
		defer close(inCopy)

		go func() {
			for v := range inCopy {
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}()
	}()
	return out
}
