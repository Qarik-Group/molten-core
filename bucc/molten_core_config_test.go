package bucc_test

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
)

func node(i int) config.NodeConfig {
	return config.NodeConfig{
		ZoneIndex: uint16(i),
		PrivateIP: net.ParseIP(fmt.Sprintf("192.168.1.%d", i+10)),
		PublicIP:  net.ParseIP(fmt.Sprintf("192.168.2.%d", i+10)),
	}
}

func cluster(s int) []config.NodeConfig {
	nodes := make([]config.NodeConfig, 0)
	for i := 0; i < s; i++ {
		nodes = append(nodes, node(i))
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
	return nodes
}

var _ = Describe("MoltenCoreConfig", func() {
	Context("given a 1 node cluster", func() {
		It("renders config", func() {
			nodes := cluster(1)
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "public_ips": { "z0": "192.168.2.10" },
  "scaling": {
    "odd3": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z0" ], "instances": 1 },
      "slice3": { "azs": [ "z0" ], "instances": 1 }
    },
    "odd5": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z0" ], "instances": 1 }
    },
    "max1": {
      "slice1": { "azs": [ "z0" ], "instances": 1 }
    },
    "max2": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z0" ], "instances": 1 },
      "slice3": { "azs": [ "z0" ], "instances": 1 },
      "slice4": { "azs": [ "z0" ], "instances": 1 },
      "slice5": { "azs": [ "z0" ], "instances": 1 }
    },
    "max3": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z0" ], "instances": 1 },
      "slice3": { "azs": [ "z0" ], "instances": 1 }
    },
    "all": {
      "x1": { "azs": [ "z0" ], "instances": 1 },
      "x2": { "azs": [ "z0" ], "instances": 2 },
      "x4": { "azs": [ "z0" ], "instances": 4 },
      "x8": { "azs": [ "z0" ], "instances": 8 },
      "x16": { "azs": [ "z0" ], "instances": 16 },
      "x32": { "azs": [ "z0" ], "instances": 32 },
      "x64": { "azs": [ "z0" ], "instances": 64 }
    }
  }
}
`))
		})
	})

	Context("given a 2 node cluster", func() {
		It("renders config", func() {
			nodes := cluster(2)
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "public_ips": {
    "z0": "192.168.2.10",
    "z1": "192.168.2.11"
  },
  "scaling": {
    "odd3": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z1" ], "instances": 1 },
      "slice3": { "azs": [ "z0" ], "instances": 1 }
    },
    "odd5": {
      "slice1": { "azs": [ "z0" ], "instances": 1 },
      "slice2": { "azs": [ "z1" ], "instances": 1 }
    },
    "max1": {
      "slice1": { "azs": [ "z0" ], "instances": 1 }
    },
    "max2": {
      "slice1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice2": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice3": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice4": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice5": { "azs": [ "z0", "z1" ], "instances": 2 }
    },
    "max3": {
      "slice1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice2": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice3": { "azs": [ "z0", "z1" ], "instances": 2 }
    },
    "all": {
      "x1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "x2": { "azs": [ "z0", "z1" ], "instances": 4 },
      "x4": { "azs": [ "z0", "z1" ], "instances": 8 },
      "x8": { "azs": [ "z0", "z1" ], "instances": 16 },
      "x16": { "azs": [ "z0", "z1" ], "instances": 32 },
      "x32": { "azs": [ "z0", "z1" ], "instances": 64 },
      "x64": { "azs": [ "z0", "z1" ], "instances": 128 }
    }
  }
}
`))
		})
	})

	Context("given a 3 node cluster", func() {
		It("renders config", func() {
			nodes := cluster(3)
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "public_ips": {
    "z0": "192.168.2.10",
    "z1": "192.168.2.11",
    "z2": "192.168.2.12"
  },
  "scaling": {
    "odd3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice3": { "azs": [ "z0", "z1", "z2" ], "instances": 3 }
    },
    "odd5": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z0", "z1", "z2" ], "instances": 3 }
    },
    "max1": {
      "slice1": { "azs": [ "z0" ], "instances": 1 }
    },
    "max2": {
      "slice1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice2": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice3": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice4": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice5": { "azs": [ "z0", "z1" ], "instances": 2 }
    },
    "max3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice3": { "azs": [ "z0", "z1", "z2" ], "instances": 3 }
    },
    "all": {
      "x1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "x2": { "azs": [ "z0", "z1", "z2" ], "instances": 6 },
      "x4": { "azs": [ "z0", "z1", "z2" ], "instances": 12 },
      "x8": { "azs": [ "z0", "z1", "z2" ], "instances": 24 },
      "x16": { "azs": [ "z0", "z1", "z2" ], "instances": 48 },
      "x32": { "azs": [ "z0", "z1", "z2" ], "instances": 96 },
      "x64": { "azs": [ "z0", "z1", "z2" ], "instances": 192 }
    }
  }
}
`))
		})
	})

	Context("given a 5 node cluster", func() {
		It("renders config", func() {
			nodes := cluster(5)
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "public_ips": {
    "z0": "192.168.2.10",
    "z1": "192.168.2.11",
    "z2": "192.168.2.12",
    "z3": "192.168.2.13",
    "z4": "192.168.2.14"
  },
  "scaling": {
    "odd3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice3": { "azs": [ "z0", "z1", "z2" ], "instances": 3 }
    },
    "odd5": {
      "slice1": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 5 },
      "slice2": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 5 }
    },
    "max1": {
      "slice1": { "azs": [ "z0" ], "instances": 1 }
    },
    "max2": {
      "slice1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice2": { "azs": [ "z2", "z3" ], "instances": 2 },
      "slice3": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice4": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice5": { "azs": [ "z0", "z1" ], "instances": 2 }
    },
    "max3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice3": { "azs": [ "z0", "z1", "z2" ], "instances": 3 }
    },
    "all": {
      "x1": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 5 },
      "x2": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 10 },
      "x4": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 20 },
      "x8": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 40 },
      "x16": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 80 },
      "x32": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 160 },
      "x64": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 320 }
    }
  }
}
`))
		})
	})

	Context("given a 9 node cluster", func() {
		It("renders config", func() {
			nodes := cluster(9)
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "public_ips": {
    "z0": "192.168.2.10",
    "z1": "192.168.2.11",
    "z2": "192.168.2.12",
    "z3": "192.168.2.13",
    "z4": "192.168.2.14",
    "z5": "192.168.2.15",
    "z6": "192.168.2.16",
    "z7": "192.168.2.17",
    "z8": "192.168.2.18"
  },
  "scaling": {
    "odd3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z3", "z4", "z5" ], "instances": 3 },
      "slice3": { "azs": [ "z6", "z7", "z8" ], "instances": 3 }
    },
    "odd5": {
      "slice1": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 5 },
      "slice2": { "azs": [ "z0", "z1", "z2", "z3", "z4" ], "instances": 5 }
    },
    "max1": {
      "slice1": { "azs": [ "z0" ], "instances": 1 }
    },
    "max2": {
      "slice1": { "azs": [ "z0", "z1" ], "instances": 2 },
      "slice2": { "azs": [ "z2", "z3" ], "instances": 2 },
      "slice3": { "azs": [ "z4", "z5" ], "instances": 2 },
      "slice4": { "azs": [ "z6", "z7" ], "instances": 2 },
      "slice5": { "azs": [ "z0", "z1" ], "instances": 2 }
    },
    "max3": {
      "slice1": { "azs": [ "z0", "z1", "z2" ], "instances": 3 },
      "slice2": { "azs": [ "z3", "z4", "z5" ], "instances": 3 },
      "slice3": { "azs": [ "z6", "z7", "z8" ], "instances": 3 }
    },
    "all": {
      "x1": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 9 },
      "x2": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 18 },
      "x4": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 36 },
      "x8": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 72 },
      "x16": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 144 },
      "x32": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 288 },
      "x64": { "azs": [ "z0", "z1", "z2", "z3", "z4", "z5", "z6", "z7", "z8" ], "instances": 576 }
    }
  }
}
`))
		})
	})

})
