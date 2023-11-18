package s3router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func dummy() *string {
	return aws.String(uuid.NewString())
}

var testCases = []struct {
	fn func(*s3.Client) error

	expected Action
}{
	{
		fn: func(client *s3.Client) error {
			_, err := client.AbortMultipartUpload(
				context.TODO(),
				&s3.AbortMultipartUploadInput{
					Bucket:   dummy(),
					Key:      dummy(),
					UploadId: dummy(),
				},
			)
			return err
		},
		expected: ActionAbortMultipartUpload,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.CompleteMultipartUpload(
				context.TODO(),
				&s3.CompleteMultipartUploadInput{
					Bucket:   dummy(),
					Key:      dummy(),
					UploadId: dummy(),
				},
			)
			return err
		},
		expected: ActionCompleteMultipartUpload,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.CopyObject(
				context.TODO(),
				&s3.CopyObjectInput{
					Bucket:     dummy(),
					Key:        dummy(),
					CopySource: dummy(),
				},
			)
			return err
		},
		expected: ActionCopyObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.CreateBucket(
				context.TODO(),
				&s3.CreateBucketInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionCreateBucket,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.CreateMultipartUpload(
				context.TODO(),
				&s3.CreateMultipartUploadInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionCreateMultipartUpload,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucket(
				context.TODO(),
				&s3.DeleteBucketInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucket,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketAnalyticsConfiguration(
				context.TODO(),
				&s3.DeleteBucketAnalyticsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketAnalyticsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketCors(
				context.TODO(),
				&s3.DeleteBucketCorsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketCors,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketEncryption(
				context.TODO(),
				&s3.DeleteBucketEncryptionInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketEncryption,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketIntelligentTieringConfiguration(
				context.TODO(),
				&s3.DeleteBucketIntelligentTieringConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketIntelligentTieringConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketInventoryConfiguration(
				context.TODO(),
				&s3.DeleteBucketInventoryConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketInventoryConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketLifecycle(
				context.TODO(),
				&s3.DeleteBucketLifecycleInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketLifecycle,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketMetricsConfiguration(
				context.TODO(),
				&s3.DeleteBucketMetricsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketMetricsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketOwnershipControls(
				context.TODO(),
				&s3.DeleteBucketOwnershipControlsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketOwnershipControls,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketPolicy(
				context.TODO(),
				&s3.DeleteBucketPolicyInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketPolicy,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketReplication(
				context.TODO(),
				&s3.DeleteBucketReplicationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketReplication,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketTagging(
				context.TODO(),
				&s3.DeleteBucketTaggingInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteBucketWebsite(
				context.TODO(),
				&s3.DeleteBucketWebsiteInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteBucketWebsite,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteObject(
				context.TODO(),
				&s3.DeleteObjectInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteObjects(
				context.TODO(),
				&s3.DeleteObjectsInput{
					Bucket: dummy(),
					Delete: &types.Delete{
						Objects: []types.ObjectIdentifier{{Key: dummy()}},
					},
				},
			)
			return err
		},
		expected: ActionDeleteObjects,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeleteObjectTagging(
				context.TODO(),
				&s3.DeleteObjectTaggingInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionDeleteObjectTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.DeletePublicAccessBlock(
				context.TODO(),
				&s3.DeletePublicAccessBlockInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionDeletePublicAccessBlock,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketAccelerateConfiguration(
				context.TODO(),
				&s3.GetBucketAccelerateConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketAccelerateConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketAcl(
				context.TODO(),
				&s3.GetBucketAclInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketACL,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketAnalyticsConfiguration(
				context.TODO(),
				&s3.GetBucketAnalyticsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketAnalyticsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketCors(
				context.TODO(),
				&s3.GetBucketCorsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketCors,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketEncryption(
				context.TODO(),
				&s3.GetBucketEncryptionInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketEncryption,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketIntelligentTieringConfiguration(
				context.TODO(),
				&s3.GetBucketIntelligentTieringConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketIntelligentTieringConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketInventoryConfiguration(
				context.TODO(),
				&s3.GetBucketInventoryConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketInventoryConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketLifecycleConfiguration(
				context.TODO(),
				&s3.GetBucketLifecycleConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketLifecycleConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketLocation(
				context.TODO(),
				&s3.GetBucketLocationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketLocation,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketLogging(
				context.TODO(),
				&s3.GetBucketLoggingInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketLogging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketMetricsConfiguration(
				context.TODO(),
				&s3.GetBucketMetricsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketMetricsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketNotificationConfiguration(
				context.TODO(),
				&s3.GetBucketNotificationConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketNotificationConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketOwnershipControls(
				context.TODO(),
				&s3.GetBucketOwnershipControlsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketOwnershipControls,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketPolicy(
				context.TODO(),
				&s3.GetBucketPolicyInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketPolicy,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketPolicyStatus(
				context.TODO(),
				&s3.GetBucketPolicyStatusInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketPolicyStatus,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketReplication(
				context.TODO(),
				&s3.GetBucketReplicationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketReplication,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketRequestPayment(
				context.TODO(),
				&s3.GetBucketRequestPaymentInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketRequestPayment,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketTagging(
				context.TODO(),
				&s3.GetBucketTaggingInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketVersioning(
				context.TODO(),
				&s3.GetBucketVersioningInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketVersioning,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetBucketWebsite(
				context.TODO(),
				&s3.GetBucketWebsiteInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetBucketWebsite,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObject(
				context.TODO(),
				&s3.GetObjectInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectAcl(
				context.TODO(),
				&s3.GetObjectAclInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectACL,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectAttributes(
				context.TODO(),
				&s3.GetObjectAttributesInput{
					Bucket: dummy(),
					Key:    dummy(),
					ObjectAttributes: []types.ObjectAttributes{
						types.ObjectAttributesEtag,
					},
				},
			)
			return err
		},
		expected: ActionGetObjectAttributes,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectLegalHold(
				context.TODO(),
				&s3.GetObjectLegalHoldInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectLegalHold,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectLockConfiguration(
				context.TODO(),
				&s3.GetObjectLockConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectLockConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectRetention(
				context.TODO(),
				&s3.GetObjectRetentionInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectRetention,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectTagging(
				context.TODO(),
				&s3.GetObjectTaggingInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetObjectTorrent(
				context.TODO(),
				&s3.GetObjectTorrentInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionGetObjectTorrent,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.GetPublicAccessBlock(
				context.TODO(),
				&s3.GetPublicAccessBlockInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionGetPublicAccessBlock,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.HeadBucket(
				context.TODO(),
				&s3.HeadBucketInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionHeadBucket,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.HeadObject(
				context.TODO(),
				&s3.HeadObjectInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionHeadObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListBucketAnalyticsConfigurations(
				context.TODO(),
				&s3.ListBucketAnalyticsConfigurationsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListBucketAnalyticsConfigurations,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListBucketIntelligentTieringConfigurations(
				context.TODO(),
				&s3.ListBucketIntelligentTieringConfigurationsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListBucketIntelligentTieringConfigurations,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListBucketInventoryConfigurations(
				context.TODO(),
				&s3.ListBucketInventoryConfigurationsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListBucketInventoryConfigurations,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListBucketMetricsConfigurations(
				context.TODO(),
				&s3.ListBucketMetricsConfigurationsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListBucketMetricsConfigurations,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListBuckets(
				context.TODO(),
				&s3.ListBucketsInput{},
			)
			return err
		},
		expected: ActionListBuckets,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListMultipartUploads(
				context.TODO(),
				&s3.ListMultipartUploadsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListMultipartUploads,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListObjects(
				context.TODO(),
				&s3.ListObjectsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListObjects,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListObjectVersions(
				context.TODO(),
				&s3.ListObjectVersionsInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionListObjectVersions,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.ListParts(
				context.TODO(),
				&s3.ListPartsInput{
					Bucket:   dummy(),
					Key:      dummy(),
					UploadId: dummy(),
				},
			)
			return err
		},
		expected: ActionListParts,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketAccelerateConfiguration(
				context.TODO(),
				&s3.PutBucketAccelerateConfigurationInput{
					Bucket:                  dummy(),
					AccelerateConfiguration: &types.AccelerateConfiguration{},
				},
			)
			return err
		},
		expected: ActionPutBucketAccelerateConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketAcl(
				context.TODO(),
				&s3.PutBucketAclInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionPutBucketACL,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketAnalyticsConfiguration(
				context.TODO(),
				&s3.PutBucketAnalyticsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
					AnalyticsConfiguration: &types.AnalyticsConfiguration{
						Id:                   dummy(),
						StorageClassAnalysis: &types.StorageClassAnalysis{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketAnalyticsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketCors(
				context.TODO(),
				&s3.PutBucketCorsInput{
					Bucket: dummy(),
					CORSConfiguration: &types.CORSConfiguration{
						CORSRules: []types.CORSRule{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketCors,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketEncryption(
				context.TODO(),
				&s3.PutBucketEncryptionInput{
					Bucket: dummy(),
					ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
						Rules: []types.ServerSideEncryptionRule{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketEncryption,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketIntelligentTieringConfiguration(
				context.TODO(),
				&s3.PutBucketIntelligentTieringConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
					IntelligentTieringConfiguration: &types.IntelligentTieringConfiguration{
						Id:       dummy(),
						Status:   types.IntelligentTieringStatusDisabled,
						Tierings: []types.Tiering{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketIntelligentTieringConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketInventoryConfiguration(
				context.TODO(),
				&s3.PutBucketInventoryConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
					InventoryConfiguration: &types.InventoryConfiguration{
						Destination: &types.InventoryDestination{
							S3BucketDestination: &types.InventoryS3BucketDestination{
								Bucket: dummy(),
								Format: types.InventoryFormatCsv,
							},
						},
						Id:                     dummy(),
						IncludedObjectVersions: types.InventoryIncludedObjectVersionsAll,
						IsEnabled:              aws.Bool(true),
						Schedule: &types.InventorySchedule{
							Frequency: types.InventoryFrequencyDaily,
						},
					},
				},
			)
			fmt.Println(err)
			return err
		},
		expected: ActionPutBucketInventoryConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketLifecycleConfiguration(
				context.TODO(),
				&s3.PutBucketLifecycleConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionPutBucketLifecycleConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketLogging(
				context.TODO(),
				&s3.PutBucketLoggingInput{
					Bucket:              dummy(),
					BucketLoggingStatus: &types.BucketLoggingStatus{},
				},
			)
			return err
		},
		expected: ActionPutBucketLogging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketMetricsConfiguration(
				context.TODO(),
				&s3.PutBucketMetricsConfigurationInput{
					Bucket: dummy(),
					Id:     dummy(),
					MetricsConfiguration: &types.MetricsConfiguration{
						Id: dummy(),
					},
				},
			)
			return err
		},
		expected: ActionPutBucketMetricsConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketNotificationConfiguration(
				context.TODO(),
				&s3.PutBucketNotificationConfigurationInput{
					Bucket:                    dummy(),
					NotificationConfiguration: &types.NotificationConfiguration{},
				},
			)
			return err
		},
		expected: ActionPutBucketNotificationConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketOwnershipControls(
				context.TODO(),
				&s3.PutBucketOwnershipControlsInput{
					Bucket: dummy(),
					OwnershipControls: &types.OwnershipControls{
						Rules: []types.OwnershipControlsRule{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketOwnershipControls,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketPolicy(
				context.TODO(),
				&s3.PutBucketPolicyInput{
					Bucket: dummy(),
					Policy: dummy(),
				},
			)
			return err
		},
		expected: ActionPutBucketPolicy,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketReplication(
				context.TODO(),
				&s3.PutBucketReplicationInput{
					Bucket: dummy(),
					ReplicationConfiguration: &types.ReplicationConfiguration{
						Role:  dummy(),
						Rules: []types.ReplicationRule{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketReplication,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketRequestPayment(
				context.TODO(),
				&s3.PutBucketRequestPaymentInput{
					Bucket: dummy(),
					RequestPaymentConfiguration: &types.RequestPaymentConfiguration{
						Payer: types.PayerBucketOwner,
					},
				},
			)
			return err
		},
		expected: ActionPutBucketRequestPayment,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketTagging(
				context.TODO(),
				&s3.PutBucketTaggingInput{
					Bucket: dummy(),
					Tagging: &types.Tagging{
						TagSet: []types.Tag{},
					},
				},
			)
			return err
		},
		expected: ActionPutBucketTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketVersioning(
				context.TODO(),
				&s3.PutBucketVersioningInput{
					Bucket:                  dummy(),
					VersioningConfiguration: &types.VersioningConfiguration{},
				},
			)
			return err
		},
		expected: ActionPutBucketVersioning,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutBucketWebsite(
				context.TODO(),
				&s3.PutBucketWebsiteInput{
					Bucket:               dummy(),
					WebsiteConfiguration: &types.WebsiteConfiguration{},
				},
			)
			return err
		},
		expected: ActionPutBucketWebsite,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObject(
				context.TODO(),
				&s3.PutObjectInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionPutObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObjectAcl(
				context.TODO(),
				&s3.PutObjectAclInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionPutObjectACL,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObjectLegalHold(
				context.TODO(),
				&s3.PutObjectLegalHoldInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionPutObjectLegalHold,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObjectLockConfiguration(
				context.TODO(),
				&s3.PutObjectLockConfigurationInput{
					Bucket: dummy(),
				},
			)
			return err
		},
		expected: ActionPutObjectLockConfiguration,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObjectRetention(
				context.TODO(),
				&s3.PutObjectRetentionInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionPutObjectRetention,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutObjectTagging(
				context.TODO(),
				&s3.PutObjectTaggingInput{
					Bucket: dummy(),
					Key:    dummy(),
					Tagging: &types.Tagging{
						TagSet: []types.Tag{},
					},
				},
			)
			return err
		},
		expected: ActionPutObjectTagging,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.PutPublicAccessBlock(
				context.TODO(),
				&s3.PutPublicAccessBlockInput{
					Bucket:                         dummy(),
					PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{},
				},
			)
			return err
		},
		expected: ActionPutPublicAccessBlock,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.RestoreObject(
				context.TODO(),
				&s3.RestoreObjectInput{
					Bucket: dummy(),
					Key:    dummy(),
				},
			)
			return err
		},
		expected: ActionRestoreObject,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.SelectObjectContent(
				context.TODO(),
				&s3.SelectObjectContentInput{
					Bucket:              dummy(),
					Key:                 dummy(),
					Expression:          dummy(),
					ExpressionType:      types.ExpressionTypeSql,
					InputSerialization:  &types.InputSerialization{},
					OutputSerialization: &types.OutputSerialization{},
				},
			)
			return err
		},
		expected: ActionSelectObjectContent,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.UploadPart(
				context.TODO(),
				&s3.UploadPartInput{
					Bucket:     dummy(),
					Key:        dummy(),
					PartNumber: aws.Int32(0),
					UploadId:   dummy(),
				},
			)
			return err
		},
		expected: ActionUploadPart,
	},
	{
		fn: func(client *s3.Client) error {
			_, err := client.UploadPartCopy(
				context.TODO(),
				&s3.UploadPartCopyInput{
					Bucket:     dummy(),
					Key:        dummy(),
					CopySource: dummy(),
					PartNumber: aws.Int32(0),
					UploadId:   dummy(),
				},
			)
			return err
		},
		expected: ActionUploadPartCopy,
	},
}

type clientAsServer struct {
	AcceptedHosts []string
	T             *testing.T

	ErrorText string

	LastRoute *Route
	LastError error
}

func (c *clientAsServer) Do(r *http.Request) (*http.Response, error) {
	r.Host = r.URL.Host
	r.Header.Set("Host", r.URL.Host)

	raw, err := httputil.DumpRequest(r, false)
	require.NoError(c.T, err)
	c.T.Log(string(raw))

	c.LastRoute, c.LastError = DetermineRoute(r, c.AcceptedHosts)

	return nil, errors.New(c.ErrorText)
}

func (c *clientAsServer) reset(t *testing.T) {
	c.T = t
	c.ErrorText = uuid.NewString()
	c.LastRoute = nil
	c.LastError = nil
}

func TestDetermineAction(t *testing.T) {
	host := "s3.local-dev.example.com"
	mitm := &clientAsServer{
		AcceptedHosts: []string{host},
	}

	for _, hostnameImmutable := range []bool{true, false} {
		t.Run(fmt.Sprintf("HostnameImmutable=%t", hostnameImmutable), func(t *testing.T) {
			client := s3.NewFromConfig(aws.Config{
				Region: "local-dev",
				EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					if service != s3.ServiceID {
						return aws.Endpoint{}, fmt.Errorf("not supported service: %q", service)
					}

					return aws.Endpoint{
							URL:               "http://" + host,
							SigningRegion:     "local-dev",
							HostnameImmutable: hostnameImmutable,
						},
						nil
				}),
				Credentials:      aws.AnonymousCredentials{},
				RetryMaxAttempts: 1,
				HTTPClient:       mitm,
			})

			for _, tc := range testCases {
				t.Run(tc.expected.String(), func(t *testing.T) {
					mitm.reset(t)

					err := tc.fn(client)
					require.Error(t, err)
					require.Contains(t, err.Error(), mitm.ErrorText)

					require.NotNil(t, mitm.LastRoute)
					require.Equal(t, tc.expected, mitm.LastRoute.Action)
					require.NoError(t, mitm.LastError)
				})
			}
		})
	}
}
