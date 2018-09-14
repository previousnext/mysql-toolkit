package exec

// Pipeline comprised of steps to execute.
type Pipeline struct {
	steps []Step
}

// Step being executed in a pipeline.
type Step struct {
	Name string
	Func func() (Status, error)
}

// Add a step function to the pipeline.
func (p *Pipeline) Add(name string, fn func() (Status, error)) {
	for delta, step := range p.steps {
		if step.Name == name {
			p.steps[delta] = Step{
				Name: name,
				Func: fn,
			}

			return
		}
	}

	p.steps = append(p.steps, Step{
		Name: name,
		Func: fn,
	})
}

// Run the pipeline.
func (p *Pipeline) Run(prev []StepStatus) ([]StepStatus, error) {
	updates := make([]StepStatus, len(p.steps))

	for delta, step := range p.steps {
		if val, ok := exists(step.Name, prev); ok {
			if val.Status == StatusFinished {
				updates[delta] = val

				continue
			}
		}

		updates[delta] = StepStatus{
			Name:   step.Name,
			Status: StatusRunning,
		}

		status, err := step.Func()
		if err != nil {
			updates[delta] = StepStatus{
				Name:    step.Name,
				Status:  status,
				Message: err.Error(),
			}

			return updates, err
		}

		updates[delta] = StepStatus{
			Name:   step.Name,
			Status: status,
		}
	}

	return updates, nil
}

// Helper function to check if a step exists.
func exists(name string, list []StepStatus) (StepStatus, bool) {
	for _, item := range list {
		if item.Name == name {
			return item, true
		}
	}

	return StepStatus{}, false
}
