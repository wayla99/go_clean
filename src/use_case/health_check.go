package use_case

import "context"

func (uc UseCase) HealthCheck(ctx context.Context) error {
	if err := uc.staffRepository.Health(ctx); err != nil {
		return err
	}
	return nil
}
