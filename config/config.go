package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/infani/awsgo/ssm"
	"github.com/infani/awsgo/cloudmap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/mod/semver"
)

var (
	UseCloudMap bool
	LogLevel    string
	Region      string
	Stage       string
	//Cloud Map
	Namespace string
	//AppSync
	Site   string
	ApiId  string
	ApiUrl string
	//ECS
	Cluster           string
	TaskDefinitionArn string
	ContainerName     string
	Subnet            string
	//S3
	Bucket string
	//Cloudfront
	KeyID            string
	PrivateKey       string
	CloudfrontDomain string
	DistributionId   string
	//iot
	ArchiveShadow string
	BackupShadow  string
	//OpenSearch
	OpenSearchParams openSearchParams
)

type openSearchParams struct {
	Address  string
	Account  string
	Password string
}

func init() {
	UseCloudMap, _ = strconv.ParseBool(os.Getenv("USE_CLOUD_MAP"))
	LogLevel = os.Getenv("LOG_LEVEL")
	Region = os.Getenv("AWS_REGION")
	Stage = os.Getenv("STAGE")
	Site = os.Getenv("SITE")
	Namespace = getParameter("NAMESPACE", "/vsaas/cloudmap/namespace")
	getAmplify()
	getMediaMaker()
	getDeviceIot()
	getIot()
	OpenSearchParams = getOpenSearchParams()
}

func getParameter(env string, name string) string {
	param := os.Getenv(env)
	if param == "" {
		var err error
		param, err = ssm.GetParameter(aws.String(name))
		if err != nil {
			return ""
		}
	}
	return param
}

func getAmplify() {
	if UseCloudMap {
		instance, err := getService("amplify")
		if err != nil {
			log.Println(err)
			return
		}
		ApiId = instance["appsyncID"]
		ApiUrl = instance["appsyncUrl"]
	} else {
		ApiId = os.Getenv("API_ID")
		ApiUrl = os.Getenv("API_URL")
	}
}

func getMediaMaker() {
	if UseCloudMap {
		instance, err := getService("mediaMaker")
		if err != nil {
			log.Println(err)
			return
		}
		Cluster = instance["cluster"]
		TaskDefinitionArn = instance["taskDefinitionArn"]
		ContainerName = instance["containerName"]
		Subnet = instance["subnet"]
	} else {
		Cluster = os.Getenv("CLUSTER")
		TaskDefinitionArn = os.Getenv("TASK_DEFINITION_ARN")
		ContainerName = os.Getenv("CONTAINER_NAME")
		Subnet = os.Getenv("SUBNET")
	}
}

func getDeviceIot() {
	if UseCloudMap {
		instance, err := getService("deviceIot")
		if err != nil {
			log.Println(err)
			return
		}
		KeyID = instance["cloudfrontKeyID"]
		CloudfrontDomain = instance["cloudfrontDomain"]
		DistributionId = instance["DistributionId"]
		Bucket = instance["s3BucketName"]
		cloudfrontKeyParamterName := instance["cloudfrontKeyParamterName"]
		PrivateKey = getParameter("PRIVATE_KEY", cloudfrontKeyParamterName)
	} else {
		Bucket = os.Getenv("BUCKET")
		KeyID = os.Getenv("KEY_ID")
		CloudfrontDomain = os.Getenv("CLOUDFRONT_DOMAIN")
		DistributionId = os.Getenv("DISTRIBUTION_ID")
	}
}

func getIot() {
	if UseCloudMap {
		ArchiveShadow = "reco"
		BackupShadow = "backup"
	} else {
		ArchiveShadow = os.Getenv("ARCHIVE_SHADOW")
		BackupShadow = os.Getenv("BACKUP_SHADOW")
	}
}

func getOpenSearchParams() openSearchParams {
	instances, err := cloudmap.DiscoverInstances("opensearch", "opensearch")
	if err != nil {
		return openSearchParams{}
	}
	if len(instances) > 0 {
		return openSearchParams{
			Address:  instances[0].Attributes["address"],
			Account:  instances[0].Attributes["account"],
			Password: instances[0].Attributes["password"],
		}
	}
	return openSearchParams{}
}

func getService(service string) (map[string]string, error) {
	instances, err := cloudmap.DiscoverInstances(Namespace, service)
	if err != nil {
		return nil, err
	}

	lastVersion := "v0.0.0"
	var result map[string]string = nil
	for _, instance := range instances {
		if instance.Attributes["site"] == Site && instance.Attributes["stage"] == Stage {
			version := instance.Attributes["version"]
			version = "v" + version
			if semver.Compare(lastVersion, version) == -1 {
				lastVersion = version
				result = instance.Attributes
			}
			// log.Println(instance.Attributes)
		}
	}
	if result == nil {
		return nil, fmt.Errorf("service : %s not found", service)
	}
	return result, nil
}

func SetByGoTest() {
	Namespace = "vivoreco"
	Region = "ap-northeast-1"
	Stage = "dev"
	Site = "site"
	instance, err := getService("vivoreco")
	if err != nil {
		log.Println(err)
		return
	}
	Bucket = instance["s3BucketName"]
	CloudfrontDomain = instance["cloudfrontDomain"]
	KeyID = instance["cloudfrontKeyID"]
	cloudfrontKeyParamterName := instance["cloudfrontKeyParamterName"]
	PrivateKey = getParameter("PRIVATE_KEY", cloudfrontKeyParamterName)
}

func toString() {
	fmt.Println("Region: ", Region)
	fmt.Println("Stage: ", Stage)
	fmt.Println("Namespace: ", Namespace)
	fmt.Println("Site : ", Site)
	fmt.Println("ApiId: ", ApiId)
	fmt.Println("ApiUrl: ", ApiUrl)
	fmt.Println("Cluster: ", Cluster)
	fmt.Println("TaskDefinitionArn: ", TaskDefinitionArn)
	fmt.Println("ContainerName: ", ContainerName)
	fmt.Println("Subnet: ", Subnet)
	fmt.Println("Bucket: ", Bucket)
	fmt.Println("KeyID: ", KeyID)
	fmt.Println("PrivateKey: ", PrivateKey)
	fmt.Println("CloudfrontDomain: ", CloudfrontDomain)
	fmt.Println("ArchiveShadow: ", ArchiveShadow)
	fmt.Println("BackupShadow: ", BackupShadow)
	fmt.Println("DistributionId: ", DistributionId)
	fmt.Println("OpenSearchParams: ", OpenSearchParams)
}
