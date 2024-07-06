package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		// go func() {
		// 	<-done
		// 	close(out) 
		// }()

		OutChan := in

		for _, stage := range stages {
			OutChan = stage(OutChan)
		}

		go func() {
			for v := range OutChan {
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
