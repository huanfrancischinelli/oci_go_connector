package functions

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/database"
	"github.com/oracle/oci-go-sdk/filestorage"
	"github.com/oracle/oci-go-sdk/identity"
	"github.com/oracle/oci-go-sdk/loadbalancer"
	"github.com/oracle/oci-go-sdk/monitoring"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"github.com/oracle/oci-go-sdk/usageapi"
)

func ListAllInstances(configProvider common.ConfigurationProvider, tenancyID string) ([]core.Instance, error) {
	var instances []core.Instance

	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating OCI client: %v", err)
		return nil, err
	}

	request := core.ListInstancesRequest{
		CompartmentId: common.String(tenancyID),
	}

	for {
		response, err := computeClient.ListInstances(context.Background(), request)
		if err != nil {
			return nil, err
		}

		instances = append(instances, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return instances, nil
}

func ListAllBootVolumes(configProvider common.ConfigurationProvider, tenancyID string, availabilityDomain string) ([]core.BootVolume, error) {
	var bootVolumes []core.BootVolume

	computeClient, err := core.NewBlockstorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating OCI client: %v", err)
		return nil, err
	}

	request := core.ListBootVolumesRequest{
		CompartmentId:      common.String(tenancyID),
		AvailabilityDomain: common.String(availabilityDomain),
	}

	for {
		response, err := computeClient.ListBootVolumes(context.Background(), request)
		if err != nil {
			return nil, err
		}

		bootVolumes = append(bootVolumes, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return bootVolumes, nil
}

func ListAllLoadBalancers(configProvider common.ConfigurationProvider, tenancyID string) ([]loadbalancer.LoadBalancer, error) {
	var loadBalancers []loadbalancer.LoadBalancer

	loadBalancerClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Load Balancer client: %v", err)
		return nil, err
	}

	request := loadbalancer.ListLoadBalancersRequest{
		CompartmentId: common.String(tenancyID),
	}

	for {
		response, err := loadBalancerClient.ListLoadBalancers(context.Background(), request)
		if err != nil {
			return nil, err
		}

		loadBalancers = append(loadBalancers, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return loadBalancers, nil
}

func ListAllBuckets(configProvider common.ConfigurationProvider, tenancyID string) ([]objectstorage.BucketSummary, error) {
	var buckets []objectstorage.BucketSummary

	objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Object Storage client: %v", err)
		return nil, err
	}
	namespaceRequest := objectstorage.GetNamespaceRequest{
		CompartmentId: common.String(tenancyID),
	}
	namespace, _ := objectStorageClient.GetNamespace(context.Background(), namespaceRequest)
	request := objectstorage.ListBucketsRequest{
		NamespaceName: common.String(*namespace.Value),
		CompartmentId: common.String(tenancyID),
	}

	for {
		response, err := objectStorageClient.ListBuckets(context.Background(), request)
		if err != nil {
			return nil, err
		}

		buckets = append(buckets, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return buckets, nil
}

func ListAvailabilityDomains(configProvider common.ConfigurationProvider, tenancyID string) ([]identity.AvailabilityDomain, error) {
	var availabilityDomains []identity.AvailabilityDomain

	identityClient, err := identity.NewIdentityClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Object Storage client: %v", err)
		return nil, err
	}

	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: common.String(tenancyID),
	}

	for {
		response, err := identityClient.ListAvailabilityDomains(context.Background(), request)
		if err != nil {
			return nil, err
		}

		availabilityDomains = append(availabilityDomains, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		// request.Page = response.OpcNextPage
	}

	return availabilityDomains, nil
}

func ListAllFileSystems(configProvider common.ConfigurationProvider, tenancyID string, availabilityDomain string) ([]filestorage.FileSystemSummary, error) {
	var fileSystems []filestorage.FileSystemSummary

	fileStorageClient, err := filestorage.NewFileStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Object Storage client: %v", err)
		return nil, err
	}

	request := filestorage.ListFileSystemsRequest{
		CompartmentId:      common.String(tenancyID),
		AvailabilityDomain: common.String(availabilityDomain),
	}

	for {
		response, err := fileStorageClient.ListFileSystems(context.Background(), request)
		if err != nil {
			return nil, err
		}

		fileSystems = append(fileSystems, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return fileSystems, nil
}

func ListFileStorageSnapshots(configProvider common.ConfigurationProvider, fileSystemID string) ([]filestorage.SnapshotSummary, error) {
	var snapshots []filestorage.SnapshotSummary

	fileStorageClient, err := filestorage.NewFileStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Object Storage client: %v", err)
		return nil, err
	}

	request := filestorage.ListSnapshotsRequest{
		FileSystemId: common.String(fileSystemID),
	}

	for {
		response, err := fileStorageClient.ListSnapshots(context.Background(), request)
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return snapshots, nil
}

func ListAllDatabases(configProvider common.ConfigurationProvider, tenancyID string) ([]database.DbSystemSummary, error) {
	var databases []database.DbSystemSummary

	databaseClient, err := database.NewDatabaseClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Object Storage client: %v", err)
		return nil, err
	}

	request := database.ListDbSystemsRequest{
		CompartmentId: common.String(tenancyID),
	}

	for {
		response, err := databaseClient.ListDbSystems(context.Background(), request)
		if err != nil {
			return nil, err
		}

		databases = append(databases, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return databases, nil
}

func GetBillingReport(configProvider common.ConfigurationProvider, tenancyID string, startDate time.Time, endDate time.Time) ([]usageapi.UsageSummary, error) {
	var billingReports []usageapi.UsageSummary

	usageapiClient, err := usageapi.NewUsageapiClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating Usage API client: %v", err)
		return nil, err
	}

	request := usageapi.RequestSummarizedUsagesRequest{

		RequestSummarizedUsagesDetails: usageapi.RequestSummarizedUsagesDetails{
			TenantId:         common.String(tenancyID),
			TimeUsageStarted: &common.SDKTime{Time: startDate},
			TimeUsageEnded:   &common.SDKTime{Time: endDate},
			Granularity:      usageapi.RequestSummarizedUsagesDetailsGranularityDaily,
		},
	}

	for {
		response, err := usageapiClient.RequestSummarizedUsages(context.Background(), request)
		if err != nil {
			return nil, err
		}

		billingReports = append(billingReports, response.Items...)

		if response.OpcNextPage == nil {
			break
		}

		request.Page = response.OpcNextPage
	}

	return billingReports, nil
}

func GetInstanceMetrics(configProvider common.ConfigurationProvider, tenancyID string, instanceID string) (map[string]float64, error) {

	monitoringClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error creating monitoring client: %v", err)
	}

	metrics := make(map[string]float64)
	metricQueries := make(map[string]string)
	metricQueries["CpuUtilization"] = fmt.Sprintf(`CpuUtilization[1m]{resourceID =~ "%s"}.max()`, instanceID)
	metricQueries["MemoryUtilization"] = fmt.Sprintf(`MemoryUtilization[1m]{resourceID =~ "%s"}.max()`, instanceID)
	metricQueries["DiskBytesRead"] = fmt.Sprintf(`DiskBytesRead[1m]{resourceID =~ "%s"}.max()`, instanceID)
	metricQueries["DiskBytesWritten"] = fmt.Sprintf(`DiskBytesWritten[1m]{resourceID =~ "%s"}.max()`, instanceID)
	for metric, metricQuery := range metricQueries {
		metrics[metric] = 0
		request := monitoring.SummarizeMetricsDataRequest{
			CompartmentId: common.String(tenancyID),
			SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
				Namespace: common.String("oci_computeagent"),
				Query:     common.String(metricQuery),
				StartTime: &common.SDKTime{Time: time.Now().Add(-1 * time.Minute)},
			},
		}

		response, err := monitoringClient.SummarizeMetricsData(context.Background(), request)
		if err != nil {
			return nil, err
		}

		for _, metricValues := range response.Items {
			lastIndex := len(metricValues.AggregatedDatapoints) - 1
			if *metricValues.AggregatedDatapoints[lastIndex].Value != 0 {
				metrics[metric] = *metricValues.AggregatedDatapoints[lastIndex].Value
			}
		}
	}
	return metrics, nil
}
