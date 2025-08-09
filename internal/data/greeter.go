package data

import (
	"context"
	"sync"

	"pet-angel/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// GreeterRepo provides a trivial in-memory implementation for demo greeter.
type GreeterRepo struct {
	data   *Data
	logger *log.Helper

	mu     sync.RWMutex
	store  map[int64]*biz.Greeter
	nextID int64
}

func NewGreeterRepo(d *Data, logger log.Logger) *GreeterRepo {
	return &GreeterRepo{data: d, logger: log.NewHelper(logger), store: make(map[int64]*biz.Greeter), nextID: 1}
}

func (r *GreeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	r.nextID++
	// No id field on biz.Greeter; keep simple
	r.store[id] = &biz.Greeter{Hello: g.Hello}
	return g, nil
}

func (r *GreeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *GreeterRepo) FindByID(ctx context.Context, id int64) (*biz.Greeter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.store[id]; ok {
		return v, nil
	}
	return &biz.Greeter{Hello: "World"}, nil
}

func (r *GreeterRepo) ListByHello(ctx context.Context, hello string) ([]*biz.Greeter, error) {
	return []*biz.Greeter{{Hello: hello}}, nil
}

func (r *GreeterRepo) ListAll(ctx context.Context) ([]*biz.Greeter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*biz.Greeter, 0, len(r.store))
	for _, v := range r.store {
		result = append(result, v)
	}
	return result, nil
}
