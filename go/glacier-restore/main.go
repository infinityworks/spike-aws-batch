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
	pageHandler := func(output *s3.ListObjectsOutput, lastPage bool) (ok bool) {
		for _, itm := range output.Contents {
			if *itm.StorageClass != glacierClass {
				continue
			}
			keysChan <- *itm.Key
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
var dryRunFlag = flag.Bool("dryRun", true, "set to false to actually run the command")
var verboseFlag = flag.Bool("verbose", true, "set verbose output")

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
	if *verboseFlag {
		log.Printf("daysFlag: %d\n tierFlag: %s\n regionFlag: %s\n bucketFlag: %s\n concurrencyFlag: %d\n dryRunFlag: %t\n verboseFlag: %t\n", *daysFlag, *tierFlag, *regionFlag, *bucketFlag, *concurrencyFlag, *dryRunFlag, *verboseFlag)
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
				if *verboseFlag {
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
