package cmd

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/urfave/cli/v2"
	"log"
)

func getConsulClient(addr, token string) (*consul.Client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	config.Token = token

	return consul.NewClient(config)
}

func CloneCommand(c *cli.Context) error {
	if !c.Bool("services") && !c.Bool("kv") {
		return fmt.Errorf("must specify at least one of --services or --kv")
	}

	source := c.String("source")
	dest := c.String("dest")

	if len(source) == 0 || len(dest) == 0 {
		return fmt.Errorf("must specify both --source and --dest")
	}

	sourceClient, err := getConsulClient(source, c.String("source-token"))
	if err != nil {
		return fmt.Errorf("failed to create source client: %v", err)
	}

	destClient, err := getConsulClient(dest, c.String("dest-token"))
	if err != nil {
		return fmt.Errorf("failed to create destination client: %v", err)
	}

	if c.Bool("kv") {
		if err := CloneKVStore(sourceClient, destClient, nil); err != nil {
			return fmt.Errorf("failed to clone KV store: %v", err)
		}
	}

	if c.Bool("services") {
		if err := CloneServices(sourceClient, destClient, nil); err != nil {
			return fmt.Errorf("failed to clone services: %v", err)
		}
	}

	return nil
}

func CloneKVStore(source, dest *consul.Client, targets []string) error {
	kv := source.KV()
	destKV := dest.KV()
	clonedCount := 0

	// If no targets specified, clone everything
	if len(targets) == 0 {
		targets = []string{""}
	}

	for _, target := range targets {
		pairs, _, err := kv.List(target, nil)
		if err != nil {
			return fmt.Errorf("failed to list KV pairs for prefix %s: %v", target, err)
		}

		for _, pair := range pairs {
			_, err := destKV.Put(&consul.KVPair{
				Key:   pair.Key,
				Value: pair.Value,
			}, nil)
			if err != nil {
				return fmt.Errorf("failed to put KV pair %s: %v", pair.Key, err)
			}
			clonedCount++
		}
	}

	log.Printf("Successfully cloned %d KV pairs", clonedCount)
	return nil
}

func CloneServices(source, dest *consul.Client, targets []string) error {
	catalog := source.Catalog()
	services, _, err := catalog.Services(nil)
	if err != nil {
		return err
	}

	// Convert targets to a map for easier lookup
	targetMap := make(map[string]bool)
	if len(targets) > 0 {
		for _, t := range targets {
			targetMap[t] = true
		}
	}

	clonedCount := 0
	for serviceName := range services {
		// Skip consul service itself
		if serviceName == "consul" {
			continue
		}

		// Skip if not in targets (when targets are specified)
		if len(targetMap) > 0 && !targetMap[serviceName] {
			continue
		}

		nodes, _, err := catalog.Service(serviceName, "", nil)
		if err != nil {
			return fmt.Errorf("failed to get service %s: %v", serviceName, err)
		}

		for _, node := range nodes {
			//registration := &consul.CatalogRegistration{
			//	Node:     node.Node,
			//	Address:  node.ServiceAddress,
			//	NodeMeta: node.NodeMeta,
			//	Service: &consul.AgentService{
			//		ID:      node.ServiceID,
			//		Service: node.ServiceName,
			//		Tags:    node.ServiceTags,
			//		Port:    node.ServicePort,
			//	},
			//	Check: nil, // Checks need to be re-created on the destination
			//}

			registration := &consul.CatalogRegistration{
				//ID:              node,
				Node:            node.Node,
				Address:         node.Address,
				TaggedAddresses: node.TaggedAddresses,
				NodeMeta:        node.NodeMeta,
				Datacenter:      node.Datacenter,
				Service: &consul.AgentService{
					//Kind:              nod,
					ID:      node.ServiceID,
					Service: node.ServiceName,
					Tags:    node.ServiceTags,
					//Meta:    node.NodeMeta,
					Port:    node.ServicePort,
					Address: node.ServiceAddress,
				},
				//Check:           node.,
				Checks:         node.Checks,
				SkipNodeUpdate: true,
				Partition:      node.Partition,
				Locality:       node.ServiceLocality,
			}
			_, err := dest.Catalog().Register(registration, nil)
			if err != nil {
				return fmt.Errorf("failed to register service %s: %v", serviceName, err)
			}
		}
		clonedCount++
		log.Printf("Cloned service: %s", serviceName)
	}

	log.Printf("Successfully cloned %d services", clonedCount)
	return nil
}

func DeregisterCommand(c *cli.Context) error {
	client, err := getConsulClient(c.String("dest"), c.String("dest-token"))
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		return fmt.Errorf("failed to list services: %v", err)
	}

	svcName := c.String("name")
	for serviceName := range services {
		// Skip consul service itself
		if serviceName == "consul" {
			continue
		}

		if len(svcName) > 0 && serviceName != svcName {
			continue
		}

		nodes, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			return fmt.Errorf("failed to get service %s: %v", serviceName, err)
		}

		for _, node := range nodes {
			deregistration := &consul.CatalogDeregistration{
				Node:      node.Node,
				ServiceID: node.ServiceID,
			}

			_, err := client.Catalog().Deregister(deregistration, nil)
			if err != nil {
				return fmt.Errorf("failed to deregister service %s: %v", serviceName, err)
			}
		}
	}

	log.Printf("Successfully deregistered all services")
	return nil
}
