package conf

import "github.com/kelseyhightower/envconfig"

var Instance *conf

type dbConf struct {
	Driver string `envconfig:"DB_DRIVER" required:"true"`
	DSN    string `envconfig:"DB_DSN" required:"true"`
}

type MinioConf struct {
	Endpoint  string `envconfig:"MINIO_ENDPOINT" required:"true"`
	AccessKey string `envconfig:"MINIO_ACCESS_KEY" required:"true"`
	SecretKey string `envconfig:"MINIO_SECRET_KEY" required:"true"`
	UseSSL    bool   `envconfig:"MINIO_USE_SSL" required:"true"`
}

type Kafka struct {
	Host        string `envconfig:"KAFKA_HOST" required:"true"`
	FileTopic   string `envconfig:"FILE_TOPIC" required:"true"`
	EventTopic  string `envconfig:"EVENT_TOPIC" required:"true"`
	FinishTopic string `envconfig:"FINISH_TOPIC" required:"true"`
}

type conf struct {
	Port     string `envconfig:"PORT" required:"true"`
	Database dbConf
	Minio    MinioConf
	Kafka    Kafka
}

func init() {
	newConf()
}

func newConf() {

	s := conf{}

	err := envconfig.Process("", &s)

	if err != nil {
		panic(err)
	}

	Instance = &s

}
