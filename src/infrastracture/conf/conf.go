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
	Host       string `envconfig:"KAFKA_HOST" required:"true"`
	FileTopic  string `envconfig:"FILE_TOPIC" required:"true"`
	EventTopic string `envconfig:"EVENT_TOPIC" required:"true"`
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

	// portenv := os.Getenv("PORT")

	// if portenv == "" {
	// 	fmt.Println("PORT not defined, using default 8080")

	// 	portenv = "8080"
	// }

	// dbdriver := os.Getenv("DB_DRIVER")

	// if dbdriver == "" {
	// 	fmt.Println("DB_DRIVER not defined, using default mysql")

	// 	dbdriver = "mysql"
	// }

	// dbdsn := os.Getenv("DB_DSN")

	// if dbdsn == "" {
	// 	fmt.Println("DB_DSN not defined... panic")

	// 	panic("DB_DSN not defined")
	// }

	// minioEndpoint := os.Getenv("MINIO_ENDPOINT")

	// if minioEndpoint == "" {
	// 	panic("MINIO_ENDPOINT not defined")
	// }

	// minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")

	// if minioAccessKey == "" {
	// 	panic("MINIO_ACCESS_KEY not defined")
	// }

	// minioSecretKey := os.Getenv("MINIO_SECRET_KEY")

	// if minioSecretKey == "" {
	// 	panic("MINIO_SECRET_KEY not defined")
	// }

	// minioUseSSL := os.Getenv("MINIO_USE_SSL")

	// if minioUseSSL == "" {
	// 	panic("MINIO_USE_SSL not defined")
	// }

	// kafkahost := os.Getenv("KAFKA_HOST")

	// if kafkahost == "" {
	// 	panic("KAFKA_HOST not defined")
	// }

	// fileTopic := os.Getenv("FILE_TOPIC")

	// if fileTopic == "" {
	// 	panic("FILE_TOPIC not defined")
	// }

	// eventTopic := os.Getenv("EVENT_TOPIC")

	// if eventTopic == "" {
	// 	panic("EVENT_TOPIC not defined")
	// }

	// Instance = &conf{
	// 	Port: portenv,
	// 	Database: dbConf{
	// 		Driver: dbdriver,
	// 		DSN:    dbdsn,
	// 	},
	// 	Minio: MinioConf{
	// 		Endpoint:  minioEndpoint,
	// 		AccessKey: minioAccessKey,
	// 		SecretKey: minioSecretKey,
	// 		UseSSL:    minioUseSSL == "true",
	// 	},
	// 	Kafka: Kafka{
	// 		Host: kafkahost,
	// 	},
	// }
}
