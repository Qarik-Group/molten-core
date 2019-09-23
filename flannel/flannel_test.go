package flannel_test

import (
	"encoding/json"

	. "github.com/starkandwayne/molten-core/flannel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subnet", func() {
	var subnet Subnet

	It("(Un)Marshals to/from JSON", func() {
		raw := []byte(`"10.1.4.0/24"`)
		Expect(json.Unmarshal(raw, &subnet)).ToNot(HaveOccurred())
		out, err := json.Marshal(subnet)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(Equal(string(raw)))
	})
})
