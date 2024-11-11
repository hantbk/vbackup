package resticProxy

import "context"

// RunLoadIndex runs the LoadIndex method on the repository.
func RunLoadIndex(repoid int) error {
	repoHandler, err := GetRepository(repoid)
	if err != nil {
		return err
	}
	repo := repoHandler.repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = repo.LoadIndex(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
