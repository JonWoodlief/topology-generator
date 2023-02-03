package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"gopkg.in/yaml.v2"
)

func createGraphDiagram(regions []string, connections []map[string]string, title string) {
	g := graphviz.New()
	graph, _ := g.Graph()

	regionNodeMap := make(map[string]*cgraph.Node)
	for _, region := range regions {
		node, _ := graph.CreateNode(region)
		regionNodeMap[region] = node
	}

	i := 0
	for _, connection := range connections {
		for node1, node2 := range connection {
			graph.CreateEdge(fmt.Sprint(i), regionNodeMap[node1], regionNodeMap[node2])
			i += 1
		}
	}

	if err := g.RenderFilename(graph, graphviz.PNG, title); err != nil {
		log.Fatal(err)
	}
}

func createRingTopology(regions []string) []map[string]string {
	topology := make([]map[string]string, 0)

	// create connections between regions
	for i, region := range regions {
		var nextRegion string
		if i == len(regions)-1 {
			nextRegion = regions[0]
		} else {
			nextRegion = regions[i+1]
		}

		topology = append(topology, map[string]string{region: nextRegion})
	}

	return topology
}

func createBidirectionalRingTopology(regions []string) []map[string]string {
	topology := make([]map[string]string, 0)

	for i, region := range regions {
		var nextRegion string
		if i == len(regions)-1 {
			nextRegion = regions[0]
		} else {
			nextRegion = regions[i+1]
		}

		topology = append(topology, map[string]string{region: nextRegion})
		topology = append(topology, map[string]string{nextRegion: region})
	}

	return topology
}
func createConnectors(regions []string) []map[string]string {
	topology := make([]map[string]string, 0)

	for i := 3; i < len(regions); i += 4 {
		oppositeRegion := regions[(i+len(regions)/2)%len(regions)]
		fmt.Printf("region %s is directly across from node %s in the ring\n", regions[i], oppositeRegion)
		topology = append(topology, map[string]string{regions[i]: oppositeRegion})
		topology = append(topology, map[string]string{oppositeRegion: regions[i]})
	}

	return topology
}

func printConnectionYaml(connections []map[string]string, title string) {
	output, err := yaml.Marshal(connections)
	if err != nil {
		log.Fatalf("error marshaling output: %v", err)
	}
	fmt.Println(title)
	fmt.Println(string(output))
}

func main() {
	data, err := ioutil.ReadFile("regions.yaml")
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	var regions []string
	err = yaml.Unmarshal(data, &regions)
	if err != nil {
		fmt.Println("Failed to parse YAML:", err)
		return
	}

	fmt.Println("Parsed Regions:", regions)

	ringTopology := createRingTopology(regions)
	fmt.Println(ringTopology)
	printConnectionYaml(ringTopology, "ring topology")
	createGraphDiagram(regions, ringTopology, "images/ringtopo.png")

	bidirectionalRingTopology := createBidirectionalRingTopology(regions)
	printConnectionYaml(bidirectionalRingTopology, "bidirectional ring topology")
	createGraphDiagram(regions, bidirectionalRingTopology, "images/biring.png")

	connectors := createConnectors(regions)
	printConnectionYaml(connectors, "connectors")
	createGraphDiagram(regions, connectors, "images/connectors.png")

	ringTopologyWithConnectors := append(connectors, ringTopology...)
	printConnectionYaml(ringTopologyWithConnectors, "ring topology with connectors")
	createGraphDiagram(regions, ringTopologyWithConnectors, "images/ringtopoconnectors.png")

	bidirectionalRingTopologyWithConnectors := append(connectors, bidirectionalRingTopology...)
	printConnectionYaml(bidirectionalRingTopologyWithConnectors, "bidirectional ring topology with connectors")
	createGraphDiagram(regions, bidirectionalRingTopologyWithConnectors, "images/biringtopoconnectors.png")
}
