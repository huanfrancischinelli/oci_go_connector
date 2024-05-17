package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"oci_go_connector/functions"

	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/common"
)

func main() {

	// Procura o arquivo .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
		os.Exit(1)
	}

	// Obtém o diretório de trabalho atual
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error fetching working dir: %v", err)
		os.Exit(1)
	}

	// Cria a string contendo o path do private_key a partir da variavel OCI_CONFIG_private_key_filename
	// fornecida no arquivo .env e salva o resultado como a variavel de ambiente OCI_CONFIG_private_key_path
	// que será usada como parametro de configuracao do SDK
	private_key_path := filepath.Join(currentDir, "/config/", os.Getenv("OCI_CONFIG_private_key_filename"))
	os.Setenv("OCI_CONFIG_private_key_path", private_key_path)

	// Gerar o configProvider a partir das variaveis de ambiente
	configProvider := common.ConfigurationProviderEnvironmentVariables("OCI_CONFIG", "")
	isConfigValid, _ := common.IsConfigurationProviderValid(configProvider)
	if isConfigValid {
		// INSTANCES
		fmt.Println("Instances:")

		instanceMetrics := make(map[string]map[string]float64)
		tenancyID := os.Getenv("OCI_CONFIG_tenancy_ocid")
		instances, err := functions.ListAllInstances(configProvider, tenancyID)
		if err != nil {
			log.Fatalf("Error fetching instances: %v", err)
			os.Exit(1)
		}

		for _, instance := range instances {
			// fmt.Println(instance)
			fmt.Printf("  Name: %s, Shape: %s, State: %s\n", *instance.DisplayName, *instance.Shape, instance.LifecycleState)

			fmt.Printf("    Shape Details:\n")
			fmt.Printf("      Processor: %s\n", *instance.ShapeConfig.ProcessorDescription)
			fmt.Printf("      Cores: %.0f\n", *instance.ShapeConfig.Ocpus)
			fmt.Printf("      Memory: %.0f Gb\n", *instance.ShapeConfig.MemoryInGBs)
			if instance.LifecycleState == "RUNNING" {
				instanceMetrics[*instance.Id], err = functions.GetInstanceMetrics(configProvider, tenancyID, *instance.Id)
				if err != nil {
					log.Fatalf("Error fetching instance %s metrics: %v", *instance.DisplayName, err)
					os.Exit(1)
				}
				fmt.Printf("    Metrics:\n")
				fmt.Printf("      CpuUtilization: %.2f %%\n", instanceMetrics[*instance.Id]["CpuUtilization"])
				fmt.Printf("      MemoryUtilization: %.2f %%\n", instanceMetrics[*instance.Id]["MemoryUtilization"])
				fmt.Printf("      DiskBytesRead: %f Mb\n", instanceMetrics[*instance.Id]["DiskBytesRead"]/(1024*1024))
				fmt.Printf("      DiskBytesWritten: %f Mb\n", instanceMetrics[*instance.Id]["DiskBytesWritten"]/(1024*1024))
			}
		}

		// LOAD BALANCERS
		fmt.Println("Load Balancers:")
		loadBalancers, err := functions.ListAllLoadBalancers(configProvider, tenancyID)
		if err != nil {
			log.Fatalf("Error fetching Load Balancers: %v", err)
			os.Exit(1)
		}

		for _, loadBalancer := range loadBalancers {
			fmt.Println(loadBalancer)
			fmt.Printf("  Name: %s, Shape: %s, State: %s\n", *loadBalancer.DisplayName, *loadBalancer.ShapeName, loadBalancer.LifecycleState)

		}

		// DATABASES
		fmt.Println("Databases:")
		databases, err := functions.ListAllDatabases(configProvider, tenancyID)
		if err != nil {
			log.Fatalf("Error fetching Databases: %v", err)
			os.Exit(1)
		}

		for _, database := range databases {
			fmt.Println(database)
			fmt.Printf("  Name: %s, DB: %s, State: %s\n", *database.DisplayName, database.DatabaseEdition, database.LifecycleState)
		}

		// BUCKETS
		fmt.Println("Buckets:")
		buckets, err := functions.ListAllBuckets(configProvider, tenancyID)
		if err != nil {
			log.Fatalf("Error fetching Buckets: %v", err)
			os.Exit(1)
		}

		for _, bucket := range buckets {
			fmt.Println(bucket)
			fmt.Printf("  Name: %s, Created by: %s, Creation Date: %s\n", *bucket.Name, *bucket.CreatedBy, *bucket.TimeCreated)
		}

		// AVAILABILITY DOMAINS
		fmt.Println("Availability Domains:")
		availabilityDomains, err := functions.ListAvailabilityDomains(configProvider, tenancyID)
		if err != nil {
			log.Fatalf("Error fetching Buckets: %v", err)
			os.Exit(1)
		}

		for _, availabilityDomain := range availabilityDomains {
			fmt.Printf("  Name: %s, ID: %s\n", *availabilityDomain.Name, *availabilityDomain.Id)
		}

		// FILE SYSTEMS
		fmt.Println("File Systems:")
		for _, availabilityDomain := range availabilityDomains {
			fileSystems, err := functions.ListAllFileSystems(configProvider, tenancyID, *availabilityDomain.Name)
			if err != nil {
				log.Fatalf("Error fetching File Systems: %v", err)
				os.Exit(1)
			}

			for _, fileSystem := range fileSystems {
				fmt.Printf("  Name: %s, Size: %d Bytes, Availability Domain: %s, State: %s\n", *fileSystem.DisplayName, fileSystem.MeteredBytes, *availabilityDomain.Name, fileSystem.LifecycleState)
				snapshots, err := functions.ListFileStorageSnapshots(configProvider, *fileSystem.Id)

				if err != nil {
					log.Fatalf("Error fetching File System (%s) snapshots: %v", *fileSystem.Id, err)
					os.Exit(1)
				}
				// FILE SYSTEM SNAPSHOTS
				fmt.Println("    Snapshots:")
				for _, snapshot := range snapshots {
					fmt.Printf("  Name: %s, Creation Date: %s Bytes, State: %s\n", *snapshot.Name, snapshot.TimeCreated, snapshot.LifecycleState)
				}
			}
		}

		// BILLING
		fmt.Println("Billing:")
		billingReports, err := functions.GetBillingReport(configProvider, tenancyID, time.Date(2024, time.May, 01, 0, 0, 0, 0, time.UTC), time.Date(2024, time.May, 31, 0, 0, 0, 0, time.UTC))
		if err != nil {
			log.Fatalf("Error fetching Buckets: %v", err)
			os.Exit(1)
		}

		for _, billingReport := range billingReports {
			// fmt.Println(billingReport)
			fmt.Printf("   Date: %s, Comp. Qty: %f\n", *billingReport.TimeUsageStarted, *billingReport.ComputedQuantity)
		}

	} else {
		log.Fatalf("Error: Config file not valid")

		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println("OCI_CONFIG_private_key_path:", os.Getenv("OCI_CONFIG_private_key_path"))
		fmt.Println("OCI_CONFIG_tenancy_ocid:", os.Getenv("OCI_CONFIG_tenancy_ocid"))
		fmt.Println("OCI_CONFIG_user_ocid:", os.Getenv("OCI_CONFIG_user_ocid"))
		fmt.Println("OCI_CONFIG_fingerprint:", os.Getenv("OCI_CONFIG_fingerprint"))
		fmt.Println("OCI_CONFIG_region:", os.Getenv("OCI_CONFIG_region"))
		fmt.Println("-----------------------------------------------------------------------")

		os.Exit(1)
	}
}
