module github.com/kassisol/hbm

go 1.20

require (
	github.com/docker/docker v20.10.18+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-plugins-helpers v0.0.0-20211224144127-6eecb7beb651
	github.com/jinzhu/gorm v1.9.17-0.20200921022817-466b344ff592
	github.com/juliengk/go-docker v0.0.0-20180321163138-f5cda316edb7
	github.com/juliengk/go-mount v0.0.0-20170406141235-e7123cbaaaf6
	github.com/juliengk/go-utils v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
)

require (
	github.com/Microsoft/go-winio v0.4.8-0.20180501170546-ab35fc04b636 // indirect
	github.com/coreos/go-systemd v0.0.0-20180409111510-d1b7d058aa2a // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/go-units v0.3.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.1-0.20210111022912-b5281034e75e // indirect
	github.com/mattn/go-sqlite3 v1.14.7-0.20210218183441-ab91e9342b02 // indirect
	github.com/moby/term v0.0.0-20220808134915-39b0c02b01ae // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1.0.20180430190053-c9281466c8b2 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.3.0 // indirect
)

replace github.com/juliengk/go-utils => github.com/jonasbroms/go-utils v0.0.0-20220910194142-e74a4beb9866

replace github.com/juliengk/go-docker => github.com/kassisol/go-docker v0.0.0-20180321163138-f5cda316edb7

replace github.com/juliengk/go-mount => github.com/kassisol/go-mount v0.0.0-20170406141235-e7123cbaaaf6
