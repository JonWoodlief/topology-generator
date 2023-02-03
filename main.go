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
	g.SetLayout(graphviz.CIRCO)
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
	imageFileName := "images/" + title + ".png"
	if err := g.RenderFilename(graph, graphviz.PNG, imageFileName); err != nil {
		log.Fatal(err)
	}
	dotFileName := "dot/" + title + ".dot"
	if err := g.RenderFilename(graph, graphviz.XDOT, dotFileName); err != nil {
		log.Fatal(err)
	}
	svgFileName := "svg/" + title + ".svg"
	if err := g.RenderFilename(graph, graphviz.SVG, svgFileName); err != nil {
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

func createDirectionalConnectors(regions []string) []map[string]string {
	topology := make([]map[string]string, 0)

	for i := 1; i < (len(regions) / 2); i += 2 {
		targetIndex := (i + len(regions)/2) % len(regions)
		topology = append(topology, map[string]string{regions[i]: regions[targetIndex-1]})
		topology = append(topology, map[string]string{regions[targetIndex]: regions[i-1]})
	}

	return topology
}

func createBidirectionalConnectors(regions []string) []map[string]string {
	topology := make([]map[string]string, 0)

	for i := 1; i < (len(regions) / 2); i += 2 {
		oppositeRegion := regions[(i+len(regions)/2)%len(regions)]
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
	createGraphDiagram(regions, ringTopology, "ring-topology")

	bidirectionalRingTopology := createBidirectionalRingTopology(regions)
	printConnectionYaml(bidirectionalRingTopology, "bidirectional ring topology")
	createGraphDiagram(regions, bidirectionalRingTopology, "bidirectional-ring-topology")

	connectors := createDirectionalConnectors(regions)
	printConnectionYaml(connectors, "connectors")
	createGraphDiagram(regions, connectors, "connectors")

	ringTopologyWithConnectors := append(connectors, ringTopology...)
	printConnectionYaml(ringTopologyWithConnectors, "ring topology with connectors")
	createGraphDiagram(regions, ringTopologyWithConnectors, "ring-topology-connectors")

	bidirectionalRingTopologyWithConnectors := append(connectors, bidirectionalRingTopology...)
	printConnectionYaml(bidirectionalRingTopologyWithConnectors, "bidirectional ring topology with connectors")
	createGraphDiagram(regions, bidirectionalRingTopologyWithConnectors, "bidirectional-ring-topology-connectors")
}
