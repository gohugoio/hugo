// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

// InvalidateCloudFront invalidates the CloudFront cache for distributionID.
// It uses the default AWS credentials from the environment.
func InvalidateCloudFront(ctx context.Context, distributionID string) error {
	// SharedConfigEnable enables loading "shared config (~/.aws/config) and
	// shared credentials (~/.aws/credentials) files".
	// See https://docs.aws.amazon.com/sdk-for-go/api/aws/session/ for more
	// details.
	// This is the same codepath used by Go CDK when creating an s3 URL.
	// TODO: Update this to a Go CDK helper once available
	// (https://github.com/google/go-cloud/issues/2003).
	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return err
	}
	req := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format("20060102150405")),
			Paths: &cloudfront.Paths{
				Items:    []*string{aws.String("/*")},
				Quantity: aws.Int64(1),
			},
		},
	}
	_, err = cloudfront.New(sess).CreateInvalidationWithContext(ctx, req)
	return err
}
