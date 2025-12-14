module go.local/go

go 1.25.5

replace github.com/rug-compling/noordergraf/go => ../../go

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/rug-compling/noordergraf/go v0.0.0-20250924180634-81e89c9e21a8
	github.com/yuin/goldmark v1.7.13
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594
	go.yaml.in/yaml/v3 v3.0.4
)

require (
	github.com/dlclark/regexp2 v1.4.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)
