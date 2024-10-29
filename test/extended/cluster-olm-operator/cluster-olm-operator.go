package operatorsv1

import (
	"path/filepath"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"

	exutil "github.com/openshift/origin/test/extended/util"
	e2e "k8s.io/kubernetes/test/e2e/framework"
)

var _ = g.Describe("[sig-cluster-olm-operator] Cluster OLM Operator", func() {

	var (
		baseDir            = exutil.FixturePath("testdata", "olmv1")
		clusterCatalog     = filepath.Join(baseDir, "clustercatalog.yaml")
		clusterExtension   = filepath.Join(baseDir, "clusterExtension.yaml")
		clusterRoleBinding = filepath.Join(baseDir, "clusterrolebinding.yaml")
		serviceAccount     = filepath.Join(baseDir, "serviceaccount.yaml")
	)

	defer g.GinkgoRecover()
	oc := exutil.NewCLIWithoutNamespace("default")

	g.It("should set upgradeable true as baseline", func(ctx g.SpecContext) {
		// Check for tech preview, if this is not tech preview, bail
		if !exutil.IsTechPreviewNoUpgrade(ctx, oc.AdminConfigClient()) {
			g.Skip("Test only runs in tech-preview")
		}

		msg, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("olm", "cluster", "-o=jsonpath={range .status.conditions[*]}{.type}{' '}{.status}").Output()
		if err != nil {
			e2e.Failf("Unable to get co %s status, error:%v", msg, err)
		}
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(msg).To(o.ContainSubstring("Upgradeable True"))
	})

	g.It("should set upgradeable false if incompatible operators found", func(ctx g.SpecContext) {
		// Check for tech preview, if this is not tech preview, bail
		if !exutil.IsTechPreviewNoUpgrade(ctx, oc.AdminConfigClient()) {
			g.Skip("Test only runs in tech-preview")
		}

		err := oc.Run("apply").Args("-f", clusterRoleBinding).Execute()
		o.Expect(err).NotTo(o.HaveOccurred(), "creating cluster role binding")
		err = oc.Run("apply").Args("-f", serviceAccount).Execute()
		o.Expect(err).NotTo(o.HaveOccurred(), "creating service account")
		err = oc.Run("apply").Args("-f", clusterCatalog).Execute()
		o.Expect(err).NotTo(o.HaveOccurred(), "creating cluster catalog")
		err = oc.Run("apply").Args("-f", clusterExtension).Execute()
		o.Expect(err).NotTo(o.HaveOccurred(), "creating cluster extension")

		msg, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("olm", "cluster", "-o=jsonpath={range .status.conditions[*]}{.type}{' '}{.status}").Output()
		if err != nil {
			e2e.Failf("Unable to get co %s status, error:%v", msg, err)
		}
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(msg).To(o.ContainSubstring("Upgradeable False"))
	})
})
