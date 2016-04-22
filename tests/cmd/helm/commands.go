package helm

import (
	"github.com/deis/workflow-e2e/tests/cmd"
	"github.com/deis/workflow-e2e/tests/settings"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

// The functions in this file implement SUCCESS CASES for commonly used `git` commands.
// This allows each of these to be re-used easily in multiple contexts.

// Uninstall executes a `helm install workflow-dev` from the current directory using the provided key.
func Uninstall(chart string) {
	sess, err := cmd.Execute("helm uninstall " + chart + " -y -n deis")
	Expect(err).NotTo(HaveOccurred())
	// sess.Wait(settings.MaxEventuallyTimeout)
	// output := string(sess.Out.Contents())
	// Expect(output).To(MatchRegexp(`Done, %s:v\d deployed to Deis`, app.Name))
	Eventually(sess, settings.MaxEventuallyTimeout).Should(Exit(0))
}

func Install(chart string) {
	sess, err := cmd.Execute("helm install " + chart)
	Expect(err).NotTo(HaveOccurred())
	// sess.Wait(settings.MaxEventuallyTimeout)
	// output := string(sess.Out.Contents())
	// Expect(output).To(MatchRegexp(`Done, %s:v\d deployed to Deis`, app.Name))
	Eventually(sess, settings.MaxEventuallyTimeout).Should(Exit(0))
}
