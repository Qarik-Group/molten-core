package bucc_test

import (
	"math/rand"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
)

var _ = Describe("MoltenCoreConfig", func() {
	var (
		nodes               []config.NodeConfig
		node0, node1, node2 config.NodeConfig
	)

	BeforeEach(func() {
		node0 = config.NodeConfig{
			ZoneIndex: 0,
			PrivateIP: net.ParseIP("192.168.1.10"),
			PublicIP:  net.ParseIP("192.168.2.10"),
		}
		node1 = config.NodeConfig{
			ZoneIndex: 1,
			PrivateIP: net.ParseIP("192.168.1.11"),
			PublicIP:  net.ParseIP("192.168.2.11"),
		}
		node2 = config.NodeConfig{
			ZoneIndex: 2,
			PrivateIP: net.ParseIP("192.168.1.12"),
			PublicIP:  net.ParseIP("192.168.2.12"),
		}
	})

	Context("Given a single node cluster", func() {
		BeforeEach(func() {
			nodes = []config.NodeConfig{node0}
		})

		It("uses z0 as other_azs (so dev clusters also work)", func() {
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "singleton_az": "z0",
  "other_azs": [ "z0" ],
  "all_azs": [ "z0" ],
  "public_ips": { "z0": "192.168.2.10" },
  "sizes": {
    "other_azs": {
      "x1": 1,
      "x2": 2,
      "x4": 4,
      "x8": 8,
      "x16": 16,
      "x32": 32,
      "x64": 64
    },
    "all_azs": {
      "x1": 1,
      "x2": 2,
      "x4": 4,
      "x8": 8,
      "x16": 16,
      "x32": 32,
      "x64": 64
    }
  }
}
`))
		})
	})

	Context("Given a 3 node cluster in random order", func() {
		BeforeEach(func() {
			nodes = []config.NodeConfig{node1, node0, node2}
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
		})

		It("Renders a deterministic result", func() {
			out, err := RenderMoltenCoreConfig(&nodes)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).Should(MatchJSON(`
{
  "singleton_az": "z0",
  "other_azs": [ "z1", "z2" ],
  "all_azs": [ "z0", "z1", "z2" ],
  "public_ips": {
    "z0": "192.168.2.10",
    "z1": "192.168.2.11",
    "z2": "192.168.2.12"
  },
  "sizes": {
    "other_azs": {
      "x1": 2,
      "x2": 4,
      "x4": 8,
      "x8": 16,
      "x16": 32,
      "x32": 64,
      "x64": 128
    },
    "all_azs": {
      "x1": 3,
      "x2": 6,
      "x4": 12,
      "x8": 24,
      "x16": 48,
      "x32": 96,
      "x64": 192
    }
  }
}
`))
		})
	})
})
