package s3router

import (
	"fmt"
	"net/http"
	"net/url"
)

type routeSelector func(Route, url.Values, http.Header) Action
type routesTree map[string]map[string]routeSelector

func staticRoute(action Action) routeSelector {
	return func(Route, url.Values, http.Header) Action {
		return action
	}
}

func conditionalQueryRoute(key string, exist, notFound Action) routeSelector {
	return func(_ Route, query url.Values, _ http.Header) Action {
		if query.Has(key) {
			return exist
		}

		return notFound
	}
}

func conditionalHeaderRoute(key string, exist, notFound Action) routeSelector {
	return func(_ Route, _ url.Values, headers http.Header) Action {
		if headers.Get(key) != "" {
			return exist
		}

		return notFound
	}
}

func init() {
	for name, tree := range map[string]routesTree{
		"root":   routesForRoot,
		"bucket": routesForbucket,
		"object": routesForObject,
	} {
		if _, exists := tree[""]; !exists {
			panic(fmt.Sprintf("Routes tree for %s do not have default route", name))
		}
	}
}

var routesForRoot = routesTree{
	"": {
		http.MethodGet: staticRoute(ActionListBuckets),
	},
}

var routesForbucket = routesTree{
	"": {
		http.MethodDelete: staticRoute(ActionDeleteBucket),
		http.MethodGet:    staticRoute(ActionListObjects),
		http.MethodHead:   staticRoute(ActionHeadBucket),
		http.MethodPut:    staticRoute(ActionCreateBucket),
	},
	"accelerate": {
		http.MethodGet: staticRoute(ActionGetBucketAccelerateConfiguration),
		http.MethodPut: staticRoute(ActionPutBucketAccelerateConfiguration),
	},
	"acl": {
		http.MethodGet: staticRoute(ActionGetBucketACL),
		http.MethodPut: staticRoute(ActionPutBucketACL),
	},
	"analytics": {
		http.MethodDelete: staticRoute(ActionDeleteBucketAnalyticsConfiguration),
		http.MethodGet:    conditionalQueryRoute("id", ActionGetBucketAnalyticsConfiguration, ActionListBucketAnalyticsConfigurations),
		http.MethodPut:    staticRoute(ActionPutBucketAnalyticsConfiguration),
	},
	"cors": {
		http.MethodDelete: staticRoute(ActionDeleteBucketCors),
		http.MethodGet:    staticRoute(ActionGetBucketCors),
		http.MethodPut:    staticRoute(ActionPutBucketCors),
	},
	"delete": {
		http.MethodPost: staticRoute(ActionDeleteObjects),
	},
	"encryption": {
		http.MethodDelete: staticRoute(ActionDeleteBucketEncryption),
		http.MethodGet:    staticRoute(ActionGetBucketEncryption),
		http.MethodPut:    staticRoute(ActionPutBucketEncryption),
	},
	"intelligent-tiering": {
		http.MethodDelete: staticRoute(ActionDeleteBucketIntelligentTieringConfiguration),
		http.MethodGet:    conditionalQueryRoute("id", ActionGetBucketIntelligentTieringConfiguration, ActionListBucketIntelligentTieringConfigurations),
		http.MethodPut:    staticRoute(ActionPutBucketIntelligentTieringConfiguration),
	},
	"inventory": {
		http.MethodDelete: staticRoute(ActionDeleteBucketInventoryConfiguration),
		http.MethodGet:    conditionalQueryRoute("id", ActionGetBucketInventoryConfiguration, ActionListBucketInventoryConfigurations),
		http.MethodPut:    staticRoute(ActionPutBucketInventoryConfiguration),
	},
	"lifecycle": {
		http.MethodDelete: staticRoute(ActionDeleteBucketLifecycle),
		http.MethodGet:    staticRoute(ActionGetBucketLifecycleConfiguration),
		http.MethodPut:    staticRoute(ActionPutBucketLifecycleConfiguration),
	},
	"location": {
		http.MethodGet: staticRoute(ActionGetBucketLocation),
	},
	"logging": {
		http.MethodGet: staticRoute(ActionGetBucketLogging),
		http.MethodPut: staticRoute(ActionPutBucketLogging),
	},
	"metrics": {
		http.MethodDelete: staticRoute(ActionDeleteBucketMetricsConfiguration),
		http.MethodGet:    conditionalQueryRoute("id", ActionGetBucketMetricsConfiguration, ActionListBucketMetricsConfigurations),
		http.MethodPut:    staticRoute(ActionPutBucketMetricsConfiguration),
	},
	"notification": {
		http.MethodGet: staticRoute(ActionGetBucketNotificationConfiguration),
		http.MethodPut: staticRoute(ActionPutBucketNotificationConfiguration),
	},
	"object-lock": {
		http.MethodGet: staticRoute(ActionGetObjectLockConfiguration),
		http.MethodPut: staticRoute(ActionPutObjectLockConfiguration),
	},
	"ownershipControls": {
		http.MethodDelete: staticRoute(ActionDeleteBucketOwnershipControls),
		http.MethodGet:    staticRoute(ActionGetBucketOwnershipControls),
		http.MethodPut:    staticRoute(ActionPutBucketOwnershipControls),
	},
	"policy": {
		http.MethodDelete: staticRoute(ActionDeleteBucketPolicy),
		http.MethodGet:    staticRoute(ActionGetBucketPolicy),
		http.MethodPut:    staticRoute(ActionPutBucketPolicy),
	},
	"policyStatus": {
		http.MethodGet: staticRoute(ActionGetBucketPolicyStatus),
	},
	"publicAccessBlock": {
		http.MethodDelete: staticRoute(ActionDeletePublicAccessBlock),
		http.MethodGet:    staticRoute(ActionGetPublicAccessBlock),
		http.MethodPut:    staticRoute(ActionPutPublicAccessBlock),
	},
	"replication": {
		http.MethodDelete: staticRoute(ActionDeleteBucketReplication),
		http.MethodGet:    staticRoute(ActionGetBucketReplication),
		http.MethodPut:    staticRoute(ActionPutBucketReplication),
	},
	"requestPayment": {
		http.MethodGet: staticRoute(ActionGetBucketRequestPayment),
		http.MethodPut: staticRoute(ActionPutBucketRequestPayment),
	},
	"tagging": {
		http.MethodDelete: staticRoute(ActionDeleteBucketTagging),
		http.MethodGet:    staticRoute(ActionGetBucketTagging),
		http.MethodPut:    staticRoute(ActionPutBucketTagging),
	},
	"uploads": {
		http.MethodGet: staticRoute(ActionListMultipartUploads),
	},
	"versioning": {
		http.MethodGet: staticRoute(ActionGetBucketVersioning),
		http.MethodPut: staticRoute(ActionPutBucketVersioning),
	},
	"versions": {
		http.MethodGet: staticRoute(ActionListObjectVersions),
	},
	"website": {
		http.MethodDelete: staticRoute(ActionDeleteBucketWebsite),
		http.MethodGet:    staticRoute(ActionGetBucketWebsite),
		http.MethodPut:    staticRoute(ActionPutBucketWebsite),
	},
}

var routesForObject = routesTree{
	"": {
		http.MethodDelete: staticRoute(ActionDeleteObject),
		http.MethodGet:    staticRoute(ActionGetObject),
		http.MethodHead:   staticRoute(ActionHeadObject),
		http.MethodPut:    conditionalHeaderRoute("x-amz-copy-source", ActionCopyObject, ActionPutObject),
	},
	"acl": {
		http.MethodGet: staticRoute(ActionGetObjectACL),
		http.MethodPut: staticRoute(ActionPutObjectACL),
	},
	"attributes": {
		http.MethodGet: staticRoute(ActionGetObjectAttributes),
	},
	"legal-hold": {
		http.MethodGet: staticRoute(ActionGetObjectLegalHold),
		http.MethodPut: staticRoute(ActionPutObjectLegalHold),
	},
	"restore": {
		http.MethodPost: staticRoute(ActionRestoreObject),
	},
	"retention": {
		http.MethodGet: staticRoute(ActionGetObjectRetention),
		http.MethodPut: staticRoute(ActionPutObjectRetention),
	},
	"select": {
		http.MethodPost: staticRoute(ActionSelectObjectContent),
	},
	"tagging": {
		http.MethodDelete: staticRoute(ActionDeleteObjectTagging),
		http.MethodGet:    staticRoute(ActionGetObjectTagging),
		http.MethodPut:    staticRoute(ActionPutObjectTagging),
	},
	"torrent": {
		http.MethodGet: staticRoute(ActionGetObjectTorrent),
	},
	"uploadId": {
		http.MethodDelete: staticRoute(ActionAbortMultipartUpload),
		http.MethodGet:    staticRoute(ActionListParts),
		http.MethodPost:   staticRoute(ActionCompleteMultipartUpload),
		http.MethodPut:    conditionalHeaderRoute("x-amz-copy-source", ActionUploadPartCopy, ActionUploadPart),
	},
	"uploads": {
		http.MethodPost: staticRoute(ActionCreateMultipartUpload),
	},
}
