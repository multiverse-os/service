module github.com/multiverse-os/service

go 1.19

require (
	github.com/multiverse-os/service/pid v0.1.0
	github.com/multiverse-os/service/signal v0.1.0
)

replace (
	github.com/multiverse-os/service/pid v0.1.0 => ./pid
	github.com/multiverse-os/service/signal v0.1.0 => ./signal
)
