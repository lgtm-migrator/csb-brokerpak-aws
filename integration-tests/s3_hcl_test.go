package integration_test

import (
	"encoding/json"
	"fmt"
	testframework "github.com/cloudfoundry/cloud-service-broker/brokerpaktestframework"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	//"fmt"
	//
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("testing with plans", Label("poc"), func() {

	BeforeEach(func() {
		Expect(mockTerraform.SetTFState([]testframework.TFStateValue{})).To(Succeed())
	})

	AfterEach(func() {
		Expect(mockTerraform.Reset()).To(Succeed())
	})

	It("should provision s3", func() {
		const s3ServiceName = "csb-aws-s3-bucket"
		_, err := broker.Provision(s3ServiceName, "private", map[string]any{"bucket_name": "csb-test-bucket"})

		var outputJson tfjson.Plan


		out, err := mockTerraform.FirstTerraformInvocationPlan()
		Expect(err).ToNot(HaveOccurred())
		err = json.Unmarshal([]byte(out), &outputJson)
		Expect(err).NotTo(HaveOccurred())

		for _, change := range outputJson.ResourceChanges {
			fmt.Printf("Change Resource Name: %s Type:%s Change: %s", change.Name, change.Type, change.Change)
			Expect(change.Type).To(Equal("aws_s3_bucket"))
			Expect(change.Name).To(Equal("b"))
			Expect(change.Change.Actions).To(ConsistOf(tfjson.Action("create")))
			Expect(change.Change.After).To(MatchAllKeys(Keys{"bucket": Equal("csb-test-bucket"), "bucket_prefix": BeNil(), "force_destroy": BeFalse(), "tags": Not(BeEmpty()), "tags_all": Not(BeEmpty())}))
		}

	})

	FIt("should provision postgres", func() {
		const postgresServiceName = "csb-aws-postgresql"

		vars := []string{`-var=instance_name=`,
		`-var=cores=2`,
		`-var=db_name=sometestname`,
		`-var=labels={"label1"="value1"}`,
		`-var=storage_gb=5`,
		`-var=publicly_accessible=false`,
		`-var=multi_az=false`,
		`-var=instance_class=blah`,
		`-var=instance_name=name`,
		`-var=engine=postgres`,
		`-var=engine_version=11`,
		`-var=aws_vpc_id=vpc-0e74b4603f30e4827`,
		`-var=storage_autoscale=false`,
		`-var=storage_autoscale_limit_gb=0`,
		`-var=storage_encrypted=false`,
		`-var=parameter_group_name=nil`,
		`-var=rds_subnet_group=nil`,
		`-var=subsume=false`,
		`-var=rds_vpc_security_group_ids=`,
		`-var=allow_major_version_upgrade=false`,
		`-var=auto_minor_version_upgrade=false`,
		`-var=maintenance_window=Sun:00:00-Sun:00:00`,
		`-var=use_tls=false`,
		`-var=region=us-west-2`,
		fmt.Sprintf(`-var=aws_access_key_id=%s`, awsAccessKeyID),
		 fmt.Sprintf(`-var=aws_secret_access_key=%s`, awsSecretAccessKey)}

		plan := ShowPlan("../terraform/rds/provision",  vars...)

		for _, change := range plan.ResourceChanges {
			fmt.Printf("Change Resource Name: %s Type:%s Change: %s", change.Name, change.Type, change.Change)
			Expect(change.Type).To(Equal("aws_db_instance"))
			Expect(change.Name).To(Equal("db_instance"))
			Expect(change.Change.Actions).To(ConsistOf(tfjson.Action("create")))
			Expect(change.Change.After).To(MatchAllKeys(
				Keys{"allocated_storage": Equal("5"),
					"engine_version": Equal(11),
					"identifier": Equal("csb-test-db"),
					"instance_class": Equal("db.t3.medium"),
					"tags": Not(BeEmpty()),
					"tags_all": Not(BeEmpty())}))
		}

		// test postgres versions
		// test parameter groups
		// How fast can this run? Separate from integration tests. If it is not fast outside the broker then it is never going to be fast.

	})

})

// TODO
// how do we tie this to the broker? can we hook into what the broker executes with another terraform mock that does this under an apply?
// pros: we would get exactly what the broker would execute, including variables and workspace. We wouldn't need to do anything to the path
// cons: more convoluted
// Do we need to tie to the broker or should we just run tf?
// cons: this looks more shady
// pros: simples
//
//
func ShowPlan(dir string, args ...string) tfjson.Plan {
	fmt.Fprintf(GinkgoWriter, "Running: tf plan %s\n", strings.Join(args, " "))

	command := exec.Command("terraform", "-chdir=" + dir, "init")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
	tmpFile, _ := ioutil.TempFile("", "test-tf")
	allArgs := []string{"-chdir=" + dir, "plan", "-refresh=false",  fmt.Sprintf("-out=%s", tmpFile.Name())}
	allArgs = append(allArgs, args...)
	command = exec.Command("terraform", allArgs...)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
	command = exec.Command("terraform", "-chdir=" + dir, "show", "-json", tmpFile.Name())
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
	var outputJson tfjson.Plan
	jsonPlan, _ := ioutil.ReadAll(tmpFile)

	json.Unmarshal(jsonPlan, &outputJson)
	return outputJson
}

