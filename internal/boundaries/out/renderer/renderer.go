package renderer

import "context"

type Renderer[D any] interface {
	Render(ctx context.Context) error
	SetData(ctx context.Context, data D)
}
