module github.com/vahidmostofi/acfg

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.26
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.2.1
	github.com/joho/godotenv v1.3.0
	github.com/montanaflynn/stats v0.6.3
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/thoas/go-funk v0.7.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.20.0
