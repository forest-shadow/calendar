package jobs

type Job interface {
	Start(errCh chan error)
}

func RunJobs(jobs []Job) error {
	errCh := make(chan error)

	for _, job := range jobs {
		job.Start(errCh)
	}

	if err := <-errCh; err != nil {
		close(errCh)
		return err
	}

	return nil
}
