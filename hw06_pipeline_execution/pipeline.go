package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		stageOut := make(Bi)
		go func(stage Stage, in In, out Bi) {
			defer close(out)
			stageIn := stage(in)
			for {
				select {
				case <-done:
					return
				case v, ok := <-stageIn:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- v:
					}
				}
			}
		}(stage, out, stageOut)
		out = stageOut
	}
	return out
}
