package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

var (
	data              []string = []string{}
	nodeAssignment    []int    = []int{}
	nodeReassignment  []int    = []int{}
	angleAssignment   []int    = []int{}
	angleReassignment []int    = []int{}

	initialNodeCount = 5
	finalNodeCount   = 6
	dataCount        = 5000
)

func angleAssign(nodemultiplier, nodes, v int64) int64 {
	getAngle := float64(v % 360.0)
	if nodemultiplier == 1 {
		assignNode := getAngle / (360.0 / float64(nodes))
		if int64(assignNode) == nodes {
			return int64(assignNode - 1)
		}
		return int64(assignNode)
	}

	tempNodes := nodemultiplier * nodes
	assignTempNode := getAngle / (360.0 / float64(tempNodes))
	assignNode := int64(assignTempNode) % nodes
	if int64(assignNode) == nodes {
		return int64(assignNode - 1)
	}
	return int64(assignNode)

}

func dataBalancingCounter(assignments []int) map[string]int {
	hoho := map[string]int{}

	for _, v := range assignments {
		hoho["node"+strconv.Itoa(v)] = hoho["node"+strconv.Itoa(v)] + 1
	}
	return hoho
}

func main() {
	for i := 0; i < dataCount; i++ {
		data = append(data, "weatherinsingaporehot"+strconv.Itoa(i))
	}

	for _, v := range data {
		bi := big.NewInt(0)
		h := md5.New()
		h.Write([]byte(v))
		hexstr := hex.EncodeToString(h.Sum(nil))
		bi.SetString(hexstr, 16)

		value := bi.Int64()
		if value < 0 {
			value = value * -1
		}

		nodeAssign := value % int64(initialNodeCount)
		nodeAssignment = append(nodeAssignment, int(nodeAssign))

		nodeReassign2 := value % int64(finalNodeCount)
		nodeReassignment = append(nodeReassignment, int(nodeReassign2))

		nodeAssign3 := angleAssign(1, int64(initialNodeCount), value)
		angleAssignment = append(angleAssignment, int(nodeAssign3))

		nodeAssign4 := angleAssign(1, int64(finalNodeCount), value)
		angleReassignment = append(angleReassignment, int(nodeAssign4))

	}

	changeRequired := 0
	for i, _ := range nodeAssignment {
		if nodeAssignment[i] != nodeReassignment[i] {
			changeRequired = changeRequired + 1
		}
	}

	angleChangeRequired := 0
	for i, _ := range angleAssignment {
		if angleAssignment[i] != angleReassignment[i] {
			angleChangeRequired = angleChangeRequired + 1
		}
	}

	fmt.Printf("%v of the data is changed\n", float64(changeRequired)/float64(dataCount)*100)
	fmt.Printf("split of data:\n%v\n", dataBalancingCounter(nodeReassignment))

	fmt.Printf("%v of the data is changed\n", float64(angleChangeRequired)/float64(dataCount)*100)
	fmt.Printf("split of data for angle:\n%v\n", dataBalancingCounter(angleReassignment))

}
