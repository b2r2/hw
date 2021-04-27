package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outer := func(in, done In) Bi {
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case v, ok := <-in:
					if !ok || v == nil {
						return
					}
					select {
					case <-done:
						return
					case out <- v:
					}
				}
			}
		}()
		return out
	}

	worker := func(out, done In, stages ...Stage) Out {
		if stages[0] == nil || out == nil {
			return nil
		}
		for _, stage := range stages {
			out = stage(outer(out, done))
		}
		return out
	}

	return worker(in, done, stages...)
}
