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

//go:build !nodeploy
// +build !nodeploy

package deploy

import (
	"context"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	gcaws "gocloud.dev/aws"
)

// InvalidateCloudFront invalidates the CloudFront cache for distributionID.
// Uses AWS credentials config from the bucket URL.
func InvalidateCloudFront(ctx context.Context, target *Target) error {
	u, err := url.Parse(target.URL)
	if err != nil {
		return err
	}
	cfg, err := gcaws.V2ConfigFromURLParams(ctx, u.Query())
	if err != nil {
		return err
	}
	cf := cloudfront.NewFromConfig(cfg)
	req := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(target.CloudFrontDistributionID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format("20060102150405")),
			Paths: &types.Paths{
				Items:    []string{"/*"},
				Quantity: aws.Int32(1),
			},
		},
	}
	_, err = cf.CreateInvalidation(ctx, req)
	return err
}
