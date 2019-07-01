package cf_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("cf login runner", func() {
	var (
		err           error
		cfLoginRunner cf.LoginRunner
		scriptRunner  *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			CFDomain: "sys.test-url.io",
			OpsManager: environment.OpsManager{
				URL:      *url,
				Username: "username",
				Password: "password",
			},
		}

		cfLoginRunner = cf.LoginRunner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = cfLoginRunner.Run(data, dryRun)
	})

	It("runs the script with a cf login", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(lines).To(Equal([]string{
			`prods="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/staged/products)"`,
			`guid="$(echo "$prods" | jq -r '.[] | select(.type == "cf") | .guid')"`,
			`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"`, `user="$(echo "$creds" | jq -r .credential.value.identity)"`,
			`pass="$(echo "$creds" | jq -r .credential.value.password)"`,
			`cf login -a "api.sys.test-url.io" -u "$user" -p "$pass" --skip-ssl-validation`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("jq", "om", "cf"))
	})

	When("run with dry run set to false", func() {
		BeforeEach(func() {
			dryRun = false
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(false))
		})
	})

	When("run with dry run set to true", func() {
		BeforeEach(func() {
			dryRun = true
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(true))
		})
	})

	When("running the script succeeds", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(nil)
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
