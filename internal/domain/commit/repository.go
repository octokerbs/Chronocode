package commit

import "context"

type Repository interface {
	Store(ctx context.Context, commit *Commit) error
}
