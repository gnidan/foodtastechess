package common

type Provider func(name string, value interface{}) error
