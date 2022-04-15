package datastore

type ConfigOptions interface {
	init() (result DS, err error)
}

type DS interface {
	Put(key string, value string) (err error)
	Get(key string) (value string, err error)
	Delete(key string) (err error)
}

type LocalConfigOptions struct {
	Path string
}

func Init(opts ConfigOptions) (result DS, err error) {
	return opts.init()
}
