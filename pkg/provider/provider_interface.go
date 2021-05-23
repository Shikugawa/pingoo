package provider

import "github.com/Shikugawa/pingoo/pkg/task"

type Provider interface {
	Get() ([]task.Task, error)
}

type ProviderCallbacks interface {
	OnDelete(id string) error
}
