package datasync

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsDatasync "github.com/aws/aws-sdk-go-v2/service/datasync"
	"github.com/aws/aws-sdk-go-v2/service/datasync/types"
	"github.com/infani/awsgo/config/awsConfig"
)

type SyncS3DataInput struct {
	S3Bucket             string
	SourceDirectory      string
	DestinationDirectory string
	Files                []string
	AccessRoleArn        string
}

func genS3Location(s3Bucket string, directory string, accessRoleArn string) ([]types.LocationListEntry, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awsDatasync.NewFromConfig(cfg)
	ctx := context.Background()
	var locations []types.LocationListEntry
	createLocationS3Input := &awsDatasync.CreateLocationS3Input{
		S3BucketArn: aws.String("arn:aws:s3:::" + s3Bucket),
		S3Config: &types.S3Config{
			BucketAccessRoleArn: aws.String(accessRoleArn),
		},
		Subdirectory: aws.String(directory),
	}
	createLocationS3Output, err := client.CreateLocationS3(ctx, createLocationS3Input)
	if err != nil {
		return nil, err
	}
	location := types.LocationListEntry{
		LocationArn: createLocationS3Output.LocationArn,
		LocationUri: aws.String("s3://" + s3Bucket + "/" + directory),
	}
	locations = append(locations, location)
	return locations, err
}

func getS3Locations(s3Bucket string, directory string) ([]types.LocationListEntry, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awsDatasync.NewFromConfig(cfg)
	ctx := context.Background()

	filters := []types.LocationFilter{}
	filter := types.LocationFilter{
		Name:     types.LocationFilterNameLocationUri,
		Operator: types.OperatorEq,
		Values:   []string{"s3://" + s3Bucket + "/" + directory},
	}
	filters = append(filters, filter)
	listLocationsInput := &awsDatasync.ListLocationsInput{
		Filters: filters,
	}
	out, err := client.ListLocations(ctx, listLocationsInput)
	if err != nil {
		return nil, err
	}
	return out.Locations, err
}

func SyncS3Data(input *SyncS3DataInput) (string, error) {
	sourceLocations, err := getS3Locations(input.S3Bucket, input.SourceDirectory)
	if err != nil {
		return "", err
	}
	if len(sourceLocations) == 0 {
		sourceLocations, err = genS3Location(input.S3Bucket, input.SourceDirectory, input.AccessRoleArn)
		if err != nil {
			return "", err
		}
	}
	destinationLocations, err := getS3Locations(input.S3Bucket, input.DestinationDirectory)
	if err != nil {
		return "", err
	}
	if len(destinationLocations) == 0 {
		destinationLocations, err = genS3Location(input.S3Bucket, input.DestinationDirectory, input.AccessRoleArn)
		if err != nil {
			return "", err
		}
	}

	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awsDatasync.NewFromConfig(cfg)
	ctx := context.Background()

	includes := []types.FilterRule{}
	filterString := ""
	for _, item := range input.Files {
		if filterString != "" {
			filterString = filterString + "|"
		}
		filterString = filterString + item
	}

	filterRule := types.FilterRule{
		FilterType: types.FilterTypeSimplePattern,
		Value:      aws.String(filterString),
	}
	includes = append(includes, filterRule)
	options := &types.Options{
		TransferMode: types.TransferModeAll,
		VerifyMode:   types.VerifyModeOnlyFilesTransferred,
	}
	createTaskInput := &awsDatasync.CreateTaskInput{
		DestinationLocationArn: destinationLocations[0].LocationArn,
		SourceLocationArn:      sourceLocations[0].LocationArn,
		Options:                options,
		Includes:               includes,
	}
	createTaskOutput, err := client.CreateTask(ctx, createTaskInput)
	if err != nil {
		return "", err
	}

	startTaskExecutionInput := &awsDatasync.StartTaskExecutionInput{
		TaskArn: createTaskOutput.TaskArn,
	}
	_, err = client.StartTaskExecution(ctx, startTaskExecutionInput)
	return *createTaskOutput.TaskArn, err
}

func deleteLocation(locationArn *string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awsDatasync.NewFromConfig(cfg)
	ctx := context.Background()

	deleteLocationInput := &awsDatasync.DeleteLocationInput{
		LocationArn: locationArn,
	}
	_, err = client.DeleteLocation(ctx, deleteLocationInput)
	return err
}

func RemoveS3SyncDataTask(taskArn string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awsDatasync.NewFromConfig(cfg)
	ctx := context.Background()

	deleteTaskInput := &awsDatasync.DeleteTaskInput{
		TaskArn: aws.String(taskArn),
	}
	_, err = client.DeleteTask(ctx, deleteTaskInput)
	if err != nil {
		return err
	}
	return err
}
