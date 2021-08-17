package driver

import (
	"context"
	"fmt"

	"github.com/ory/x/configx"

	"github.com/ory/x/logrusx"

	"github.com/driver005/oauth/config"

	"github.com/driver005/oauth/registry"
)

type options struct {
	forcedValues map[string]interface{}
	preload      bool
	validate     bool
	opts         []configx.OptionModifier
}

func newOptions() *options {
	return &options{
		validate: true,
		preload:  true,
		opts:     []configx.OptionModifier{},
	}
}

type OptionsModifier func(*options)

func WithOptions(opts ...configx.OptionModifier) OptionsModifier {
	return func(o *options) {
		o.opts = append(o.opts, opts...)
	}
}

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisableValidation() OptionsModifier {
	return func(o *options) {
		o.validate = false
	}
}

// DisableValidation validating the config.
//
// This does not affect schema validation!
func DisablePreloading() OptionsModifier {
	return func(o *options) {
		o.preload = false
	}
}

func New(ctx context.Context) registry.Registry {
	o := newOptions()

	l := logrusx.New("ORY Hydra", config.Version)
	c, err := config.New(l)
	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate configuration.")
	}

	if o.validate {
		config.MustValidate(l, c)
	}

	r, err := registry.NewRegistryFromDSN(ctx, c, l)
	if err != nil {
		l.WithError(err).Fatal("Unable to create service registry.")
	}

	if err = r.Init(ctx); err != nil {
		l.WithError(err).Fatal("Unable to initialize service registry.")
	}
	fmt.Println(o.preload)
	// Avoid cold cache issues on boot:
	if o.preload {
		registry.CallRegistry(ctx, r)
	}

	c.Source().SetTracer(context.Background(), r.Tracer(ctx))

	return r
}
