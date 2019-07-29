package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const glacierClass = "GLACIER"

func getGlacierItems(ctx context.Context, region, bucket string, keysChan chan string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := s3.New(sess)
	pageHandler := func(s3objects *s3.ListObjectsOutput, lastPage bool) (ok bool) {
		for _, item := range s3objects.Contents {
			if *item.StorageClass != glacierClass {
				continue
			}
			keysChan <- *item.Key
		}
		return true
	}
	err := svc.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}, pageHandler)
	if err != nil {
		err = fmt.Errorf("failed to list items in bucket: %v", err)
		return err
	}
	return nil
}

func restoreItem(ctx context.Context, region, bucket, key, tier string, days int64) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := s3.New(sess)
	_, err := svc.RestoreObjectWithContext(ctx, &s3.RestoreObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		RestoreRequest: &s3.RestoreRequest{
			Days: aws.Int64(days),
			GlacierJobParameters: &s3.GlacierJobParameters{
				Tier: aws.String(tier),
			},
		},
	})
	return err
}

var allowedTiers = map[string]struct{}{
	"Bulk":      struct{}{},
	"Standard":  struct{}{},
	"Expedited": struct{}{},
}

var daysFlag = flag.Int("days", 1, "amount of time to keep restoration alive for")
var tierFlag = flag.String("tier", "Bulk", "glacier restore tier, affects cost, Bulk/Standard/Expedited. See https://aws.amazon.com/premiumsupport/knowledge-center/restore-glacier-tiers/")
var regionFlag = flag.String("region", "eu-west-1", "the default region to use, e.g. eu-west-1")
var bucketFlag = flag.String("bucket", "", "the glacier bucket to restore, e.g. bucket_name")
var concurrencyFlag = flag.Int("concurrency", 4, "number of concurrent aws s3 restore requests")
var dryRunFlag = flag.Bool("dryRun", false, "passed to enable a dryrun")
var quietFlag = flag.Bool("quiet", false, "passed to reduce logging output")

func areFlagsValid() bool {
	if *bucketFlag == "" {
		return false
	}
	if _, isTierAllowed := allowedTiers[*tierFlag]; !isTierAllowed {
		return false
	}
	return true
}

func main() {
	flag.Parse()
	if !*quietFlag {
    log.Printf("daysFlag: %d\n", *daysFlag);
    log.Printf("tierFlag: %s\n", *tierFlag)
    log.Printf("regionFlag: %s\n", *regionFlag)
    log.Printf("bucketFlag: %s\n", *bucketFlag)
    log.Printf("concurrencyFlag: %d\n", *concurrencyFlag)
    log.Printf("dryRunFlag: %t\n", *dryRunFlag)
    log.Printf("quietFlag: %t\n", *quietFlag)
	}
	if !areFlagsValid() {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	keysChan := make(chan string)
	var sg sync.WaitGroup
	// Start n concurrent routines.
	for i := 0; i < *concurrencyFlag; i++ {
		sg.Add(1)
		go func(channelID int) {
			defer sg.Done()
			for k := range keysChan {
				if !*dryRunFlag {
					err := restoreItem(ctx, *regionFlag, *bucketFlag, k, *tierFlag, int64(*daysFlag))
					if err != nil {
						log.Printf("error restoring %s: %v", k, err)
					}
				}
				if !*quietFlag {
					log.Printf("channel %d has processed %s", channelID, k)
				}
			}
		}(i + 1)
	}

	log.Println("Starting to retrieve glacier items.")
	err := getGlacierItems(ctx, *regionFlag, *bucketFlag, keysChan)
	if err != nil {
		fmt.Println("error getting the items:", err)
		os.Exit(1)
	}

	close(keysChan)
	sg.Wait()

	log.Print("Restoration requests completed")
}
