package influxdb

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

func InitDb(ctx context.Context, client *influxdb2.Client, orgName string, bucketName string) error {
	orgsAPI := (*client).OrganizationsAPI()
	org, err := orgsAPI.FindOrganizationByName(ctx, orgName)
	if err != nil {
		org, err = orgsAPI.CreateOrganizationWithName(ctx, orgName)
		if err != nil {
			return err
		}
	}

	bucketsAPI := (*client).BucketsAPI()
	_, err = bucketsAPI.FindBucketByName(ctx, bucketName)
	if err != nil {
		_, err = bucketsAPI.CreateBucketWithNameWithID(ctx, *org.Id, bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

func PushMeasurements(
	ctx context.Context,
	client *influxdb2.Client,
	orgName string,
	bucketName string,
	fieldName string,
	measurement string,
	tags map[string]string,
	series map[time.Time]float64,
) error {
	writeAPI := (*client).WriteAPIBlocking(orgName, bucketName)
	//allowWrite := true

	points := []*write.Point{}
	if len(series) > 0 {
		if len(series) > 1 {
			deleteAPI := (*client).DeleteAPI()
			err := deleteAPI.DeleteWithName(
				ctx,
				orgName,
				bucketName,
				time.Now().Add(-24*time.Hour),
				time.Now().Add(24*time.Hour),
				fmt.Sprintf("_measurement=\"%s\"", measurement),
			)
			if err != nil {
				return err
			}
		}

		for pointInTime, carbonIntensity := range series {
			fields := map[string]interface{}{
				fieldName: carbonIntensity,
			}
			point := write.NewPoint(measurement, tags, fields, pointInTime)
			points = append(points, point)
		}

		if err := writeAPI.WritePoint(ctx, points...); err != nil {
			return err
		}
	}

	err := writeAPI.Flush(ctx)
	if err != nil {
		return err
	}

	return nil
}
