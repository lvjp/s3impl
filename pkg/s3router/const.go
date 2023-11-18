//go:generate go run golang.org/x/tools/cmd/stringer -output const_string.go -type=Action,RequestStyle

package s3router

type Action int

const (
	ActionUnknow Action = iota
	ActionAbortMultipartUpload
	ActionCompleteMultipartUpload
	ActionCopyObject
	ActionCORSPreflightRequest
	ActionCreateBucket
	ActionCreateMultipartUpload
	ActionDeleteBucket
	ActionDeleteBucketAnalyticsConfiguration
	ActionDeleteBucketCors
	ActionDeleteBucketEncryption
	ActionDeleteBucketIntelligentTieringConfiguration
	ActionDeleteBucketInventoryConfiguration
	ActionDeleteBucketLifecycle
	ActionDeleteBucketMetricsConfiguration
	ActionDeleteBucketOwnershipControls
	ActionDeleteBucketPolicy
	ActionDeleteBucketReplication
	ActionDeleteBucketTagging
	ActionDeleteBucketWebsite
	ActionDeleteObject
	ActionDeleteObjects
	ActionDeleteObjectTagging
	ActionDeletePublicAccessBlock
	ActionGetBucketAccelerateConfiguration
	ActionGetBucketACL
	ActionGetBucketAnalyticsConfiguration
	ActionGetBucketCors
	ActionGetBucketEncryption
	ActionGetBucketIntelligentTieringConfiguration
	ActionGetBucketInventoryConfiguration
	ActionGetBucketLifecycleConfiguration
	ActionGetBucketLocation
	ActionGetBucketLogging
	ActionGetBucketMetricsConfiguration
	ActionGetBucketNotificationConfiguration
	ActionGetBucketOwnershipControls
	ActionGetBucketPolicy
	ActionGetBucketPolicyStatus
	ActionGetBucketReplication
	ActionGetBucketRequestPayment
	ActionGetBucketTagging
	ActionGetBucketVersioning
	ActionGetBucketWebsite
	ActionGetObject
	ActionGetObjectACL
	ActionGetObjectAttributes
	ActionGetObjectLegalHold
	ActionGetObjectLockConfiguration
	ActionGetObjectRetention
	ActionGetObjectTagging
	ActionGetObjectTorrent
	ActionGetPublicAccessBlock
	ActionHeadBucket
	ActionHeadObject
	ActionListBucketAnalyticsConfigurations
	ActionListBucketIntelligentTieringConfigurations
	ActionListBucketInventoryConfigurations
	ActionListBucketMetricsConfigurations
	ActionListBuckets
	ActionListMultipartUploads
	ActionListObjects
	ActionListObjectVersions
	ActionListParts
	ActionPutBucketAccelerateConfiguration
	ActionPutBucketACL
	ActionPutBucketAnalyticsConfiguration
	ActionPutBucketCors
	ActionPutBucketEncryption
	ActionPutBucketIntelligentTieringConfiguration
	ActionPutBucketInventoryConfiguration
	ActionPutBucketLifecycleConfiguration
	ActionPutBucketLogging
	ActionPutBucketMetricsConfiguration
	ActionPutBucketNotificationConfiguration
	ActionPutBucketOwnershipControls
	ActionPutBucketPolicy
	ActionPutBucketReplication
	ActionPutBucketRequestPayment
	ActionPutBucketTagging
	ActionPutBucketVersioning
	ActionPutBucketWebsite
	ActionPutObject
	ActionPutObjectACL
	ActionPutObjectLegalHold
	ActionPutObjectLockConfiguration
	ActionPutObjectRetention
	ActionPutObjectTagging
	ActionPutPublicAccessBlock
	ActionRestoreObject
	ActionSelectObjectContent
	ActionUploadPart
	ActionUploadPartCopy
)

type RequestStyle int

const (
	RequestStyleUnknown RequestStyle = iota
	RequestStyleVHost
	RequestStyleVPath
	RequestStyleCName
)

type Route struct {
	Action   Action
	Hostname string
	BaseHost string
	Style    RequestStyle
	Bucket   string
	Key      string
}
