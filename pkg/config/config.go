package config

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	LogPath     string `mapstructure:"log_path"`

	Databases map[string]Database `mapstructure:"databases"`

	Redises map[string]Redis `mapstructure:"redises"`

	OSS OSS `mapstructure:"oss"`

	OSSDirect OSS `mapstructure:"oss_direct"`

	ImageBuildService ImageBuildService `mapstructure:"image_build_service"`
	EmailAuth         EmailAuth         `mapstructure:"email_auth"`
	Auth              Auth              `mapstructure:"auth"`

	GameBaseDomain    string `mapstructure:"game_base_domain"`
	GameBackendDomain string `mapstructure:"game_backend_domain"`

	PaasTestApis PaasTestApis `mapstructure:"paas_test_apis"`

	InsideToken string `mapstructure:"inside_token"`

	BusinessDemoDomain string `mapstructure:"business_demo_domain"`

	PaasDomain string `mapstructure:"paas_domain"`

	PaasInternalDomain string `mapstructure:"paas_internal_domain"`

	Nas Nas `mapstructure:"nas"`

	Ecs Ecs `mapstructure:"ecs"`

	Kafka Kafka `mapstructure:"kafka"`

	StorageClassName string `mapstructure:"storage_class_name"`

	ImageRepository string `mapstructure:"image_repository"`

	TestBaseDomain string `mapstructure:"test_base_domain"`

	FeiShuToken string `mapstructure:"feishu_token"`
}

type Kafka struct {
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Brokers  []string `mapstructure:"brokers"`
}

type PaasTestApis struct {
	ApiGetSSToken    string `mapstructure:"api_get_sstoken"`
	ApiUpdateSSToken string `mapstructure:"api_update_sstoken"`
	ApiGetUserInfo   string `mapstructure:"api_get_userinfo"`
	ApiGetBalance    string `mapstructure:"api_get_balance"`
	ApiChangeBalance string `mapstructure:"api_change_balance"`
}

type Auth struct {
	JwtSecret string `mapstructure:"jwt_secret"`
}

type EmailAuth struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ImageBuildService struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type Redis struct {
	Host         string `mapstructure:"host"`
	InternalConn string `mapstructure:"internal_conn"`
	Port         string `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
}

type Database struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type OSS struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Domain          string `mapstructure:"domain"`
	InternalDomain  string `mapstructure:"internal_domain"`
	RoleArn         string `mapstructure:"role_arn"`
	RoleSessionName string `mapstructure:"role_session_name"`
	Region          string `mapstructure:"region"`
}

type Nas struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Endpoint        string `mapstructure:"endpoint"`
	RegionId        string `mapstructure:"region_id"`
	ZoneId          string `mapstructure:"zone_id"`
	RegionEndpoint  string `mapstructure:"region_endpoint"`
	VpcId           string `mapstructure:"vpc_id"`
	VswId           string `mapstructure:"vsw_id"`
}

type Ecs struct {
	AccessKeyId             string `mapstructure:"access_key_id"`
	AccessKeySecret         string `mapstructure:"access_key_secret"`
	Endpoint                string `mapstructure:"endpoint"`
	RegionId                string `mapstructure:"region_id"`
	ZoneId                  string `mapstructure:"zone_id"`
	VswId                   string `mapstructure:"vsw_id"`
	InstanceType            string `mapstructure:"instance_type"`
	ImageId                 string `mapstructure:"image"`
	SecurityGroupId         string `mapstructure:"security_group_id"`
	InternetChargeType      string `mapstructure:"internet_charge_type"`
	DiskCategory            string `mapstructure:"disk_category"`
	DiskSize                string `mapstructure:"disk_size"`
	InternetMaxBandwidthIn  int32  `mapstructure:"internet_max_bandwidth_in"`
	InternetMaxBandwidthOut int32  `mapstructure:"internet_max_bandwidth_out"`
	UniqueSuffix            bool   `mapstructure:"unique_suffix"`
	PasswordInherit         bool   `mapstructure:"password_inherit"`
	IoOptimized             string `mapstructure:"io_optimized"`
	Amount                  int32  `mapstructure:"amount"`
	MinAmount               int32  `mapstructure:"min_amount"`
	AutoRenew               bool   `mapstructure:"auto_renew"`
	Period                  int32  `mapstructure:"period"`
	PeriodUnit              string `mapstructure:"period_unit"`
	AutoRenewPeriod         int32  `mapstructure:"auto_renew_period"`
	InstanceChargeType      string `mapstructure:"instance_charge_type"`
}

var conf = &Config{}

func GetConfig() *Config {
	return conf
}

func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PAAS")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	bindDbEnv()
	bindRedisEnv()
	bindOssEnv()

	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.Unmarshal(conf)
	})

	viper.WatchConfig()

	return nil
}

func bindOssEnv() {
	viper.BindEnv("oss.access_key_secret")
}

func bindRedisEnv() {
	viper.BindEnv("redises.paas.password")
}

func bindDbEnv() {
	viper.BindEnv("databases.paas.password")
}

const (
	Development = "development"
	Production  = "production"
	Test        = "test"
)

func IsDevelopment() bool {
	return conf.Environment == Development
}

func IsProduction() bool {
	return conf.Environment == Production
}

const (
	LogSidecarName          = "logs-sidecar"
	LogSidecarVolumeName    = "volume-xxakk"
	LogSidecarMountPath     = "/etc/vector"
	LogSidecarConfigMapName = "volume-iefct"
	LogConfigMapName        = "vector-logs"
	LogConfigMapData        = `[sources.source_logs]
type = "file"
include = ["/app/log/*.log"]
ignore_older = 86400
read_from = "beginning"

[sinks.console]
type = "console"
inputs = ["source_logs"]
target = "stdout"
encoding.codec = "text"`
)
