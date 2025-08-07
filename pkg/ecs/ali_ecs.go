package ecs

import (
	"encoding/base64"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/usual2970/retell/pkg/config"
)

type TenantEcsInfo struct {
	AccessKeyId             string
	AccessKeySecret         string
	InstanceType            string //实例的资源规格
	InstanceName            string // 实例名称
	Endpoint                string
	RegionId                string //实例所属的地域ID,例:cn-hangzhou
	ImageId                 string //镜像ID
	SecurityGroupId         string
	InstanceChargeType      string //实际付费方式,PrePaid包年包月。
	VswId                   string //云交换机信息
	ZoneId                  string //
	UserDomain              string //用户的域名
	DiskCategory            string //系统盘类型，cloud_essd_entry为云盘
	DiskSize                string //磁盘大小，"60"单位GB
	HostName                string
	InternetMaxBandwidthIn  int32 //公网入口
	InternetMaxBandwidthOut int32 //公网出口的带宽
	UniqueSuffix            bool
	PasswordInherit         bool   //是否继承密码,镜像已经配置
	InternetChargeType      string //公网出宽带计费方式,PayByBandwidth:按固定宽带计费
	IoOptimized             string //系统盘的IO优化
	Amount                  int32  //实例数量
	MinAmount               int32  //最小实例数量
	AutoRenew               bool   //是否自动续费
	Period                  int32  //创建资源时常
	PeriodUnit              string //自动续费周期,Month为月
	AutoRenewPeriod         int32  //自动续费周期，1为一个月
	DryRun                  bool
}

func CreateEcsClient() (result *aliecs.Client, err error) {
	getConfig := config.GetConfig()
	ecsConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(getConfig.Ecs.AccessKeyId).
		SetAccessKeySecret(getConfig.Ecs.AccessKeySecret)
	akCredential, err := credentials.NewCredential(ecsConfig)
	if err != nil {
		return nil, err
	}
	openapiConfig := &openapi.Config{
		Credential: akCredential,
	}
	openapiConfig.Endpoint = tea.String(getConfig.Ecs.Endpoint)
	result = &aliecs.Client{}
	result, err = aliecs.NewClient(openapiConfig)

	return result, nil
}

func CreateEcsForTenant(client *aliecs.Client, domain, hostName string) (*aliecs.RunInstancesResponseBody, error) {
	getConfig := config.GetConfig()
	systemDisk := &aliecs.RunInstancesRequestSystemDisk{
		Size:     tea.String(getConfig.Ecs.DiskSize),
		Category: tea.String(getConfig.Ecs.DiskCategory),
	}

	createInstanceRequest := &aliecs.RunInstancesRequest{
		RegionId:                tea.String(getConfig.Ecs.RegionId),
		InternetMaxBandwidthIn:  tea.Int32(getConfig.Ecs.InternetMaxBandwidthIn),
		InternetMaxBandwidthOut: tea.Int32(getConfig.Ecs.InternetMaxBandwidthOut),
		ImageId:                 tea.String(getConfig.Ecs.ImageId),
		InstanceType:            tea.String(getConfig.Ecs.InstanceType),
		InstanceName:            tea.String(hostName),
		HostName:                tea.String(hostName),
		UniqueSuffix:            tea.Bool(getConfig.Ecs.UniqueSuffix),
		PasswordInherit:         tea.Bool(getConfig.Ecs.PasswordInherit),
		ZoneId:                  tea.String(getConfig.Ecs.ZoneId),
		InternetChargeType:      tea.String(getConfig.Ecs.InternetChargeType),
		SystemDisk:              systemDisk,
		IoOptimized:             tea.String(getConfig.Ecs.IoOptimized),
		Amount:                  tea.Int32(getConfig.Ecs.Amount),
		MinAmount:               tea.Int32(getConfig.Ecs.MinAmount),
		SecurityGroupId:         tea.String(getConfig.Ecs.SecurityGroupId),
		VSwitchId:               tea.String(getConfig.Ecs.VswId),
		PeriodUnit:              tea.String(getConfig.Ecs.PeriodUnit),
		AutoRenew:               tea.Bool(getConfig.Ecs.AutoRenew),
		AutoRenewPeriod:         tea.Int32(getConfig.Ecs.AutoRenewPeriod),
		InstanceChargeType:      tea.String(getConfig.Ecs.InstanceChargeType),
		Period:                  tea.Int32(getConfig.Ecs.Period),
		UserData:                tea.String(CreateUserScript(domain, config.GetConfig().Environment)), //启动脚本Base64编码
	}

	runtime := &util.RuntimeOptions{}

	instances, err := client.RunInstancesWithOptions(createInstanceRequest, runtime)
	if err != nil {
		return nil, err
	}
	return instances.Body, nil
}

// DescribeInstanceStatus 查询实例状态
func DescribeInstanceStatus(client *aliecs.Client, instanceIds []*string) (*aliecs.DescribeInstanceStatusResponseBody, error) {
	getConfig := config.GetConfig()
	describeInstanceStatusRequest := &aliecs.DescribeInstanceStatusRequest{
		RegionId:   tea.String(getConfig.Ecs.RegionId),
		InstanceId: instanceIds,
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.DescribeInstanceStatusWithOptions(describeInstanceStatusRequest, runtime)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func CreateUserScript(domain, env string) string {
	script := `#!/bin/bash
set -ex
app_base_domain01="%s"
cluster_evn="%s"  ##test：开发测试集群， production:生产集群


http_port="30080"
https_port="30443"
file1=/etc/motd
file2=/etc/monitrc
file3=/etc/nginx/nginx.conf
lanip=$(curl -s http://100.100.100.200/latest/meta-data/private-ipv4)
wanip=$(curl -s http://100.100.100.200/latest/meta-data/eipv4)
total_mem=$(cat /proc/meminfo | grep "MemTotal:" | awk '{print $2}')
total_mem_gb=$(echo "scale=2; $total_mem / 1024 / 1024" | bc)
cpu_cores=$(grep -c 'processor' /proc/cpuinfo)
cpu_model=$(grep 'model name' /proc/cpuinfo | head -n 1 | cut -d ':' -f 2 | sed 's/^[[:space:]]*//')

sed -i "s/CPU-INFO/${cpu_cores}x $cpu_model/g" $file1
sed -i "s/MEM-INFO/${total_mem_gb}GB/g" $file1
sed -i "s/LAN-IP/$lanip/g" "$file1" "$file2"
sed -i "s/WAN-IP/$wanip/g" "$file2"

if [[ $cluster_evn == "test" ]]; then
    ip_pool=("10.100.88.36" "10.100.88.37" "10.100.88.38")
elif [[ $cluster_evn == "production" ]]; then
    ip_pool=("10.8.0.116" "10.8.0.117" "10.8.0.118")
else
    echo "Invalid cluster environment. Please set cluster_evn to 'test' or 'production'."
fi

for ip in "${ip_pool[@]}"; do
    sed -i '/upstream 3os_http /a\        server '$ip':'$http_port' weight=1;' $file3
    sed -i '/upstream 3os_https /a\        server '$ip':'$https_port' weight=1;' $file3
done

nginx -t 
systemctl enable --now nginx monit

/usr/local/sbin/dnspod --action add --domain *.$app_base_domain01 --type A --value $wanip

rm -rf /usr/local/sbin/*`

	script = fmt.Sprintf(script, domain, env)
	// 转换为 base64
	encoded := base64.StdEncoding.EncodeToString([]byte(script))
	return encoded
}

// DescribeInstanceAttribute 查询实例属性
func DescribeInstanceAttribute(client *aliecs.Client, instanceId string) (*aliecs.DescribeInstanceAttributeResponseBody, error) {
	describeInstanceAttributeRequest := &aliecs.DescribeInstanceAttributeRequest{
		InstanceId: tea.String(instanceId),
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.DescribeInstanceAttributeWithOptions(describeInstanceAttributeRequest, runtime)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
