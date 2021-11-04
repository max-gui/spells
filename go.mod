module github.com/max-gui/spells

go 1.15

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210920160938-87db9fbc61c7 // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/bndr/gojenkins v1.1.0
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v1.8.5
	github.com/hashicorp/consul/api v1.11.0
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/max-gui/consulagent v0.0.0-20211102065914-c94b0cf85096
	github.com/max-gui/fileconvagt v0.0.0-20211102071148-98d485218484
	github.com/max-gui/logagent v0.0.0-20211102065508-44b5d1757320
	github.com/max-gui/redisagent v0.0.0-20211102070829-626b6136f5ee
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.6 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3 // indirect
	golang.org/x/sys v0.0.0-20211102061401-a2f17f7b995c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/max-gui/logagent => ../logagent

replace github.com/max-gui/consulagent => ../consulagent

replace github.com/max-gui/redisagent => ../redisagent

replace github.com/max-gui/fileconvagt => ../fileconvagt
