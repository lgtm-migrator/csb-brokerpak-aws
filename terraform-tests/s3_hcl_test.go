package terraform_tests

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gstruct"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

//var tfBinary = "/Users/mcampo/localstack/bin/tflocal"
var tfBinary = "terraform"

var _ = Describe("testing with plans", Label("poc"), func() {


	//BeforeEach(func() {
	//	Expect(mockTerraform.SetTFState([]testframework.TFStateValue{})).To(Succeed())
	//})

	//AfterEach(func() {
	//	Expect(mockTerraform.Reset()).To(Succeed())
	//})
	//
	//It("should provision s3", func() {
	//	const s3ServiceName = "csb-aws-s3-bucket"
	//	_, err := broker.Provision(s3ServiceName, "private", map[string]any{"bucket_name": "csb-test-bucket"})
	//
	//	var outputJson tfjson.Plan
	//
	//
	//	out, err := mockTerraform.FirstTerraformInvocationPlan()
	//	Expect(err).ToNot(HaveOccurred())
	//	err = json.Unmarshal([]byte(out), &outputJson)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	for _, change := range outputJson.ResourceChanges {
	//		fmt.Printf("Change Resource Name: %s Type:%s Change: %s", change.Name, change.Type, change.Change)
	//		Expect(change.Type).To(Equal("aws_s3_bucket"))
	//		Expect(change.Name).To(Equal("b"))
	//		Expect(change.Change.Actions).To(ConsistOf(tfjson.Action("create")))
	//		Expect(change.Change.After).To(MatchAllKeys(Keys{"bucket": Equal("csb-test-bucket"), "bucket_prefix": BeNil(), "force_destroy": BeFalse(), "tags": Not(BeEmpty()), "tags_all": Not(BeEmpty())}))
	//	}
	//
	//})

	BeforeEach(OncePerOrdered, func() {
		Init("../terraform/postgresql/provision")
	})

	Context("postgres parameter groups", Ordered, func() {
		const postgresServiceName = "csb-aws-postgresql"
		var plan tfjson.Plan
		var parameterGroupName string

		Context("No parameter group passed", func(){
			BeforeEach(func() {
				parameterGroupName = ""
				vars := []string{`-var=instance_name=`,
					`-var=cores=2`,
					`-var=db_name=sometestname`,
					`-var=labels={"label1"="value1"}`,
					`-var=storage_gb=5`,
					`-var=publicly_accessible=false`,
					`-var=multi_az=false`,
					`-var=instance_class=blah`,
					`-var=instance_name=name`,
					`-var=postgres_version=11`,
					`-var=aws_vpc_id=vpc-72464617`,
					`-var=storage_autoscale=false`,
					`-var=storage_autoscale_limit_gb=0`,
					`-var=storage_encrypted=false`,
					fmt.Sprintf(`-var=parameter_group_name=%s`, parameterGroupName),
					`-var=rds_subnet_group=nil`,
					`-var=subsume=false`,
					`-var=rds_vpc_security_group_ids=`,
					`-var=allow_major_version_upgrade=false`,
					`-var=auto_minor_version_upgrade=false`,
					`-var=maintenance_window=Sun:00:00-Sun:00:00`,
					`-var=region=us-west-2`,
					`-var=backup_window=00:00-00:00`,
					`-var=copy_tags_to_snapshot=true`,
					`-var=delete_automated_backups=true`,
					`-var=deletion_protection=false`,
					`-var=iops=3000`,
					`-var=kms_key_id=`,
					`-var=monitoring_interval=0`,
					`-var=monitoring_role_arn=0`,
					`-var=performance_insights_enabled=false`,
					`-var=performance_insights_kms_key_id=`,
					`-var=provider_verify_certificate=true`,
					`-var=require_ssl=false`,
					`-var=storage_type=io1`,
					`-var=backup_retention_period=7`,
					fmt.Sprintf(`-var=aws_access_key_id=%s`, awsAccessKeyID),
					fmt.Sprintf(`-var=aws_secret_access_key=%s`, awsSecretAccessKey)}
				plan = ShowPlan("../terraform/postgresql/provision", vars...)
			})

			It("should create a parameter group if no name is provided", func(){
				Expect(plan.ResourceChanges).To(HaveLen(6))

				Expect(plan.ResourceChanges).To(ContainElements(
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("aws_db_instance"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("aws_security_group_rule"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("random_string"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("random_password"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("aws_db_parameter_group"),})),

				))

			})
		})

		Context("Parameter group passed", func() {

			BeforeEach(func() {
				parameterGroupName = "some-parameter-group"
				vars := []string{`-var=instance_name=`,
					`-var=cores=2`,
					`-var=db_name=sometestname`,
					`-var=labels={"label1"="value1"}`,
					`-var=storage_gb=5`,
					`-var=publicly_accessible=false`,
					`-var=multi_az=false`,
					`-var=instance_class=blah`,
					`-var=instance_name=name`,
					`-var=postgres_version=11`,
					`-var=aws_vpc_id=vpc-72464617`,
					`-var=storage_autoscale=false`,
					`-var=storage_autoscale_limit_gb=0`,
					`-var=storage_encrypted=false`,
					fmt.Sprintf(`-var=parameter_group_name=%s`, parameterGroupName),
					`-var=rds_subnet_group=nil`,
					`-var=subsume=false`,
					`-var=rds_vpc_security_group_ids=`,
					`-var=allow_major_version_upgrade=false`,
					`-var=auto_minor_version_upgrade=false`,
					`-var=maintenance_window=Sun:00:00-Sun:00:00`,
					`-var=region=us-west-2`,
					`-var=backup_window=00:00-00:00`,
					`-var=copy_tags_to_snapshot=true`,
					`-var=delete_automated_backups=true`,
					`-var=deletion_protection=false`,
					`-var=iops=3000`,
					`-var=kms_key_id=`,
					`-var=monitoring_interval=0`,
					`-var=monitoring_role_arn=0`,
					`-var=performance_insights_enabled=false`,
					`-var=performance_insights_kms_key_id=`,
					`-var=provider_verify_certificate=true`,
					`-var=require_ssl=false`,
					`-var=storage_type=io1`,
					`-var=backup_retention_period=7`,
					fmt.Sprintf(`-var=aws_access_key_id=%s`, awsAccessKeyID),
					fmt.Sprintf(`-var=aws_secret_access_key=%s`, awsSecretAccessKey)}
				plan = ShowPlan("../terraform/postgresql/provision", vars...)
			})

			It("should not create a parameter group if name is provided", func() {
				Expect(plan.ResourceChanges).To(HaveLen(5))

				Expect(plan.ResourceChanges).To(ContainElements(
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("aws_db_instance"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("aws_security_group_rule"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("random_string"),})),
					PointTo(MatchFields(IgnoreExtras, Fields{"Type": Equal("random_password"),})),
				))
			})
		})

	})

	Context("postgres plans", Ordered, func() {
		const postgresServiceName = "csb-aws-postgresql"

		DescribeTable("should provision different postgresql versions",
			func(version, resourceQuantity int) {
				vars := []string{`-var=instance_name=`,
					`-var=cores=2`,
					`-var=db_name=sometestname`,
					`-var=labels={"label1"="value1"}`,
					`-var=storage_gb=5`,
					`-var=publicly_accessible=false`,
					`-var=multi_az=false`,
					`-var=instance_class=blah`,
					`-var=instance_name=name`,
					fmt.Sprintf(`-var=postgres_version=%d`, version),
					`-var=aws_vpc_id=vpc-72464617`,
					`-var=storage_autoscale=false`,
					`-var=storage_autoscale_limit_gb=0`,
					`-var=storage_encrypted=false`,
					`-var=parameter_group_name="`,
					`-var=rds_subnet_group=nil`,
					`-var=subsume=false`,
					`-var=rds_vpc_security_group_ids=`,
					`-var=allow_major_version_upgrade=false`,
					`-var=auto_minor_version_upgrade=false`,
					`-var=maintenance_window=Sun:00:00-Sun:00:00`,
					`-var=region=us-west-2`,
					`-var=backup_window=00:00-00:00`,
					`-var=copy_tags_to_snapshot=true`,
					`-var=delete_automated_backups=true`,
					`-var=deletion_protection=false`,
					`-var=iops=3000`,
					`-var=kms_key_id=`,
					`-var=monitoring_interval=0`,
					`-var=monitoring_role_arn=0`,
					`-var=performance_insights_enabled=false`,
					`-var=performance_insights_kms_key_id=`,
					`-var=provider_verify_certificate=true`,
					`-var=require_ssl=false`,
					`-var=storage_type=io1`,
					`-var=backup_retention_period=7`,
					fmt.Sprintf(`-var=aws_access_key_id=%s`, awsAccessKeyID),
					fmt.Sprintf(`-var=aws_secret_access_key=%s`, awsSecretAccessKey)}
				plan := ShowPlan("../terraform/postgresql/provision", vars...)

				Expect(plan.ResourceChanges).To(HaveLen(resourceQuantity))
				Expect(plan.ResourceChanges).To(ContainElement(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type": Equal("aws_db_instance"),
						"Change": PointTo(MatchFields(IgnoreExtras, Fields{
							"After": MatchKeys(IgnoreExtras, Keys{
								"engine":         Equal("postgres"),
								"engine_version": Equal(fmt.Sprint(version)),
							}),
						})),
					})),
				))
			},
			Entry("postgres 11", 11, 5),
			Entry("postgres 12", 12, 5),
			Entry("postgres 13", 13, 5),
			Entry("postgres 14", 14, 5),
		)
	})

		/// never creates anything, so don't really have a state, I could do init once and the do plans with different inputs

		//It("postgres 11", func() {
		//
		//
		//
		//
		//
		//
		//
		//
		//	// test postgres versions
		//	// test parameter groups
		//	// How fast can this run? Separate from integration tests. If it is not fast outside the broker then it is never going to be fast.
		//
		//})
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

func Init(dir string, args ...string) {
	fmt.Fprintf(GinkgoWriter, "Running: tf plan %s\n", strings.Join(args, " "))

	command := exec.Command(tfBinary, "-chdir=" + dir, "init")
	session, err := localBinaryStart(command)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
}

func ShowPlan(dir string, args ...string) tfjson.Plan {
	fmt.Fprintf(GinkgoWriter, "Running: tf plan %s\n", strings.Join(args, " "))

	tmpFile, _ := ioutil.TempFile("", "test-tf")
	allArgs := []string{"-chdir=" + dir, "plan", "-refresh=false",  fmt.Sprintf("-out=%s", tmpFile.Name())}
	allArgs = append(allArgs, args...)
	command := exec.Command("terraform", allArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
	command = exec.Command("terraform", "-chdir=" + dir, "show", "-json", tmpFile.Name())
	jsonPlan, err := command.Output()
	Expect(err).NotTo(HaveOccurred())
	var outputJson tfjson.Plan
	err = json.Unmarshal(jsonPlan, &outputJson)
	Expect(err).NotTo(HaveOccurred())
	return outputJson
}

func ShowPlan2(dir string, args ...string) tfjson.Plan {
	fmt.Fprintf(GinkgoWriter, "Running: tf plan %s\n", strings.Join(args, " "))

	tmpFile, _ := ioutil.TempFile("", "test-tf")
	allArgs := []string{"-chdir=" + dir, "plan", "-refresh=false",  fmt.Sprintf("-out=%s", tmpFile.Name())}
	allArgs = append(allArgs, args...)
	command := exec.Command(tfBinary, allArgs...)
	localBinaryStart(command)
	command = exec.Command(tfBinary, "-chdir=" + dir, "show", "-json", tmpFile.Name())
	localBinaryStart(command)
	var outputJson tfjson.Plan
	jsonPlan, _ := ioutil.ReadAll(tmpFile)

	json.Unmarshal(jsonPlan, &outputJson)
	return outputJson
}

func localBinaryStart(command *exec.Cmd) (*gexec.Session, error) {
	command.Env = append(os.Environ(),
		"VIRTUAL_ENV=/Users/mcampo/localstack",
		"PATH=/Users/mcampo/localstack/bin;" + os.Getenv("PATH"),
	)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
	return session, err

}

func loadJson() {
	json_string := "{\"format_version\":\"1.1\",\"terraform_version\":\"1.2.4\",\"variables\":{\"allow_major_version_upgrade\":{\"value\":\"false\"},\"auto_minor_version_upgrade\":{\"value\":\"false\"},\"aws_access_key_id\":{\"value\":\"\"},\"aws_secret_access_key\":{\"value\":\"\"},\"aws_vpc_id\":{\"value\":\"vpc-72464617\"},\"backup_retention_period\":{\"value\":\"7\"},\"backup_window\":{\"value\":\"00:00-00:00\"},\"copy_tags_to_snapshot\":{\"value\":\"true\"},\"cores\":{\"value\":\"2\"},\"db_name\":{\"value\":\"sometestname\"},\"delete_automated_backups\":{\"value\":\"true\"},\"deletion_protection\":{\"value\":\"false\"},\"instance_class\":{\"value\":\"blah\"},\"instance_name\":{\"value\":\"name\"},\"iops\":{\"value\":\"3000\"},\"kms_key_id\":{\"value\":\"\"},\"labels\":{\"value\":{\"label1\":\"value1\"}},\"maintenance_window\":{\"value\":\"Sun:00:00-Sun:00:00\"},\"monitoring_interval\":{\"value\":\"0\"},\"monitoring_role_arn\":{\"value\":\"0\"},\"multi_az\":{\"value\":\"false\"},\"parameter_group_name\":{\"value\":\"nil\"},\"performance_insights_enabled\":{\"value\":\"false\"},\"performance_insights_kms_key_id\":{\"value\":\"\"},\"postgres_version\":{\"value\":\"11\"},\"provider_verify_certificate\":{\"value\":\"true\"},\"publicly_accessible\":{\"value\":\"false\"},\"rds_subnet_group\":{\"value\":\"nil\"},\"rds_vpc_security_group_ids\":{\"value\":\"\"},\"region\":{\"value\":\"us-west-2\"},\"require_ssl\":{\"value\":\"false\"},\"storage_autoscale\":{\"value\":\"false\"},\"storage_autoscale_limit_gb\":{\"value\":\"0\"},\"storage_encrypted\":{\"value\":\"false\"},\"storage_gb\":{\"value\":\"5\"},\"storage_type\":{\"value\":\"io1\"},\"subsume\":{\"value\":\"false\"}},\"planned_values\":{\"outputs\":{\"hostname\":{\"sensitive\":false},\"name\":{\"sensitive\":false},\"password\":{\"sensitive\":true},\"provider_verify_certificate\":{\"sensitive\":false,\"type\":\"bool\",\"value\":true},\"require_ssl\":{\"sensitive\":false,\"type\":\"bool\",\"value\":false},\"status\":{\"sensitive\":false},\"username\":{\"sensitive\":false}},\"root_module\":{\"resources\":[{\"address\":\"aws_db_instance.db_instance\",\"mode\":\"managed\",\"type\":\"aws_db_instance\",\"name\":\"db_instance\",\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"schema_version\":1,\"values\":{\"allocated_storage\":5,\"allow_major_version_upgrade\":false,\"apply_immediately\":true,\"auto_minor_version_upgrade\":false,\"backup_retention_period\":7,\"copy_tags_to_snapshot\":true,\"customer_owned_ip_enabled\":null,\"db_name\":\"sometestname\",\"db_subnet_group_name\":\"nil\",\"delete_automated_backups\":true,\"deletion_protection\":false,\"domain\":null,\"domain_iam_role_name\":null,\"enabled_cloudwatch_logs_exports\":null,\"engine\":\"postgres\",\"engine_version\":\"11\",\"final_snapshot_identifier\":null,\"iam_database_authentication_enabled\":null,\"identifier\":\"name\",\"instance_class\":\"blah\",\"iops\":3000,\"max_allocated_storage\":null,\"monitoring_interval\":0,\"monitoring_role_arn\":\"0\",\"multi_az\":false,\"parameter_group_name\":\"nil\",\"performance_insights_enabled\":false,\"publicly_accessible\":false,\"replicate_source_db\":null,\"restore_to_point_in_time\":[],\"s3_import\":[],\"security_group_names\":null,\"skip_final_snapshot\":true,\"storage_encrypted\":false,\"storage_type\":\"io1\",\"tags\":{\"label1\":\"value1\"},\"tags_all\":{\"label1\":\"value1\"},\"timeouts\":null},\"sensitive_values\":{\"password\":true,\"replicas\":[],\"restore_to_point_in_time\":[],\"s3_import\":[],\"tags\":{},\"tags_all\":{},\"vpc_security_group_ids\":[]}},{\"address\":\"aws_security_group.rds-sg[0]\",\"mode\":\"managed\",\"type\":\"aws_security_group\",\"name\":\"rds-sg\",\"index\":0,\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"schema_version\":1,\"values\":{\"description\":\"Managed by Terraform\",\"name\":\"name-sg\",\"revoke_rules_on_delete\":false,\"tags\":null,\"timeouts\":null,\"vpc_id\":\"vpc-72464617\"},\"sensitive_values\":{\"egress\":[],\"ingress\":[],\"tags_all\":{}}},{\"address\":\"aws_security_group_rule.rds_inbound_access[0]\",\"mode\":\"managed\",\"type\":\"aws_security_group_rule\",\"name\":\"rds_inbound_access\",\"index\":0,\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"schema_version\":2,\"values\":{\"cidr_blocks\":[\"0.0.0.0/0\"],\"description\":null,\"from_port\":5432,\"ipv6_cidr_blocks\":null,\"prefix_list_ids\":null,\"protocol\":\"tcp\",\"self\":false,\"timeouts\":null,\"to_port\":5432,\"type\":\"ingress\"},\"sensitive_values\":{\"cidr_blocks\":[false]}},{\"address\":\"random_password.password\",\"mode\":\"managed\",\"type\":\"random_password\",\"name\":\"password\",\"provider_name\":\"registry.terraform.io/hashicorp/random\",\"schema_version\":2,\"values\":{\"keepers\":null,\"length\":32,\"lower\":true,\"min_lower\":0,\"min_numeric\":0,\"min_special\":0,\"min_upper\":0,\"number\":true,\"numeric\":true,\"override_special\":\"~_-.\",\"special\":false,\"upper\":true},\"sensitive_values\":{}},{\"address\":\"random_string.username\",\"mode\":\"managed\",\"type\":\"random_string\",\"name\":\"username\",\"provider_name\":\"registry.terraform.io/hashicorp/random\",\"schema_version\":2,\"values\":{\"keepers\":null,\"length\":16,\"lower\":true,\"min_lower\":0,\"min_numeric\":0,\"min_special\":0,\"min_upper\":0,\"numeric\":false,\"override_special\":null,\"special\":false,\"upper\":true},\"sensitive_values\":{}}]}},\"resource_changes\":[{\"address\":\"aws_db_instance.db_instance\",\"mode\":\"managed\",\"type\":\"aws_db_instance\",\"name\":\"db_instance\",\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"change\":{\"actions\":[\"create\"],\"before\":null,\"after\":{\"allocated_storage\":5,\"allow_major_version_upgrade\":false,\"apply_immediately\":true,\"auto_minor_version_upgrade\":false,\"backup_retention_period\":7,\"copy_tags_to_snapshot\":true,\"customer_owned_ip_enabled\":null,\"db_name\":\"sometestname\",\"db_subnet_group_name\":\"nil\",\"delete_automated_backups\":true,\"deletion_protection\":false,\"domain\":null,\"domain_iam_role_name\":null,\"enabled_cloudwatch_logs_exports\":null,\"engine\":\"postgres\",\"engine_version\":\"11\",\"final_snapshot_identifier\":null,\"iam_database_authentication_enabled\":null,\"identifier\":\"name\",\"instance_class\":\"blah\",\"iops\":3000,\"max_allocated_storage\":null,\"monitoring_interval\":0,\"monitoring_role_arn\":\"0\",\"multi_az\":false,\"parameter_group_name\":\"nil\",\"performance_insights_enabled\":false,\"publicly_accessible\":false,\"replicate_source_db\":null,\"restore_to_point_in_time\":[],\"s3_import\":[],\"security_group_names\":null,\"skip_final_snapshot\":true,\"storage_encrypted\":false,\"storage_type\":\"io1\",\"tags\":{\"label1\":\"value1\"},\"tags_all\":{\"label1\":\"value1\"},\"timeouts\":null},\"after_unknown\":{\"address\":true,\"arn\":true,\"availability_zone\":true,\"backup_window\":true,\"ca_cert_identifier\":true,\"character_set_name\":true,\"endpoint\":true,\"engine_version_actual\":true,\"hosted_zone_id\":true,\"id\":true,\"identifier_prefix\":true,\"kms_key_id\":true,\"latest_restorable_time\":true,\"license_model\":true,\"maintenance_window\":true,\"name\":true,\"nchar_character_set_name\":true,\"option_group_name\":true,\"password\":true,\"performance_insights_kms_key_id\":true,\"performance_insights_retention_period\":true,\"port\":true,\"replica_mode\":true,\"replicas\":true,\"resource_id\":true,\"restore_to_point_in_time\":[],\"s3_import\":[],\"snapshot_identifier\":true,\"status\":true,\"tags\":{},\"tags_all\":{},\"timezone\":true,\"username\":true,\"vpc_security_group_ids\":true},\"before_sensitive\":false,\"after_sensitive\":{\"password\":true,\"replicas\":[],\"restore_to_point_in_time\":[],\"s3_import\":[],\"tags\":{},\"tags_all\":{},\"vpc_security_group_ids\":[]}}},{\"address\":\"aws_security_group.rds-sg[0]\",\"mode\":\"managed\",\"type\":\"aws_security_group\",\"name\":\"rds-sg\",\"index\":0,\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"change\":{\"actions\":[\"create\"],\"before\":null,\"after\":{\"description\":\"Managed by Terraform\",\"name\":\"name-sg\",\"revoke_rules_on_delete\":false,\"tags\":null,\"timeouts\":null,\"vpc_id\":\"vpc-72464617\"},\"after_unknown\":{\"arn\":true,\"egress\":true,\"id\":true,\"ingress\":true,\"name_prefix\":true,\"owner_id\":true,\"tags_all\":true},\"before_sensitive\":false,\"after_sensitive\":{\"egress\":[],\"ingress\":[],\"tags_all\":{}}}},{\"address\":\"aws_security_group_rule.rds_inbound_access[0]\",\"mode\":\"managed\",\"type\":\"aws_security_group_rule\",\"name\":\"rds_inbound_access\",\"index\":0,\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"change\":{\"actions\":[\"create\"],\"before\":null,\"after\":{\"cidr_blocks\":[\"0.0.0.0/0\"],\"description\":null,\"from_port\":5432,\"ipv6_cidr_blocks\":null,\"prefix_list_ids\":null,\"protocol\":\"tcp\",\"self\":false,\"timeouts\":null,\"to_port\":5432,\"type\":\"ingress\"},\"after_unknown\":{\"cidr_blocks\":[false],\"id\":true,\"security_group_id\":true,\"source_security_group_id\":true},\"before_sensitive\":false,\"after_sensitive\":{\"cidr_blocks\":[false]}}},{\"address\":\"random_password.password\",\"mode\":\"managed\",\"type\":\"random_password\",\"name\":\"password\",\"provider_name\":\"registry.terraform.io/hashicorp/random\",\"change\":{\"actions\":[\"create\"],\"before\":null,\"after\":{\"keepers\":null,\"length\":32,\"lower\":true,\"min_lower\":0,\"min_numeric\":0,\"min_special\":0,\"min_upper\":0,\"number\":true,\"numeric\":true,\"override_special\":\"~_-.\",\"special\":false,\"upper\":true},\"after_unknown\":{\"bcrypt_hash\":true,\"id\":true,\"result\":true},\"before_sensitive\":false,\"after_sensitive\":{\"bcrypt_hash\":true,\"result\":true}}},{\"address\":\"random_string.username\",\"mode\":\"managed\",\"type\":\"random_string\",\"name\":\"username\",\"provider_name\":\"registry.terraform.io/hashicorp/random\",\"change\":{\"actions\":[\"create\"],\"before\":null,\"after\":{\"keepers\":null,\"length\":16,\"lower\":true,\"min_lower\":0,\"min_numeric\":0,\"min_special\":0,\"min_upper\":0,\"numeric\":false,\"override_special\":null,\"special\":false,\"upper\":true},\"after_unknown\":{\"id\":true,\"number\":true,\"result\":true},\"before_sensitive\":false,\"after_sensitive\":{}}}],\"output_changes\":{\"hostname\":{\"actions\":[\"create\"],\"before\":null,\"after_unknown\":true,\"before_sensitive\":false,\"after_sensitive\":false},\"name\":{\"actions\":[\"create\"],\"before\":null,\"after_unknown\":true,\"before_sensitive\":false,\"after_sensitive\":false},\"password\":{\"actions\":[\"create\"],\"before\":null,\"after_unknown\":true,\"before_sensitive\":true,\"after_sensitive\":true},\"provider_verify_certificate\":{\"actions\":[\"create\"],\"before\":null,\"after\":true,\"after_unknown\":false,\"before_sensitive\":false,\"after_sensitive\":false},\"require_ssl\":{\"actions\":[\"create\"],\"before\":null,\"after\":false,\"after_unknown\":false,\"before_sensitive\":false,\"after_sensitive\":false},\"status\":{\"actions\":[\"create\"],\"before\":null,\"after_unknown\":true,\"before_sensitive\":false,\"after_sensitive\":false},\"username\":{\"actions\":[\"create\"],\"before\":null,\"after_unknown\":true,\"before_sensitive\":false,\"after_sensitive\":false}},\"prior_state\":{\"format_version\":\"1.0\",\"terraform_version\":\"1.2.4\",\"values\":{\"outputs\":{\"provider_verify_certificate\":{\"sensitive\":false,\"value\":true,\"type\":\"bool\"},\"require_ssl\":{\"sensitive\":false,\"value\":false,\"type\":\"bool\"}},\"root_module\":{\"resources\":[{\"address\":\"data.aws_subnets.all\",\"mode\":\"data\",\"type\":\"aws_subnets\",\"name\":\"all\",\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"schema_version\":0,\"values\":{\"filter\":[{\"name\":\"vpc-id\",\"values\":[\"vpc-72464617\"]}],\"id\":\"us-west-2\",\"ids\":[\"subnet-0f401bc14d50a1c1c\",\"subnet-91e9a8c8\"],\"tags\":null,\"timeouts\":null},\"sensitive_values\":{\"filter\":[{\"values\":[false]}],\"ids\":[false,false]}},{\"address\":\"data.aws_vpc.vpc\",\"mode\":\"data\",\"type\":\"aws_vpc\",\"name\":\"vpc\",\"provider_name\":\"registry.terraform.io/hashicorp/aws\",\"schema_version\":0,\"values\":{\"arn\":\"arn:aws:ec2:us-west-2:649758297924:vpc/vpc-72464617\",\"cidr_block\":\"172.31.0.0/16\",\"cidr_block_associations\":[{\"association_id\":\"vpc-cidr-assoc-ea127b83\",\"cidr_block\":\"172.31.0.0/16\",\"state\":\"associated\"}],\"default\":true,\"dhcp_options_id\":\"dopt-886ec0ed\",\"enable_dns_hostnames\":true,\"enable_dns_support\":true,\"filter\":null,\"id\":\"vpc-72464617\",\"instance_tenancy\":\"default\",\"ipv6_association_id\":\"\",\"ipv6_cidr_block\":\"\",\"main_route_table_id\":\"rtb-842b36e1\",\"owner_id\":\"649758297924\",\"state\":null,\"tags\":{},\"timeouts\":null},\"sensitive_values\":{\"cidr_block_associations\":[{}],\"tags\":{}}}]}}},\"configuration\":{\"provider_config\":{\"aws\":{\"name\":\"aws\",\"full_name\":\"registry.terraform.io/hashicorp/aws\",\"version_constraint\":\"\\u003e= 4.0.0\",\"expressions\":{\"access_key\":{\"references\":[\"var.aws_access_key_id\"]},\"region\":{\"references\":[\"var.region\"]},\"secret_key\":{\"references\":[\"var.aws_secret_access_key\"]}}},\"random\":{\"name\":\"random\",\"full_name\":\"registry.terraform.io/hashicorp/random\",\"version_constraint\":\"\\u003e= 3.3.2\"}},\"root_module\":{\"outputs\":{\"hostname\":{\"expression\":{\"references\":[\"aws_db_instance.db_instance.address\",\"aws_db_instance.db_instance\"]}},\"name\":{\"expression\":{\"references\":[\"aws_db_instance.db_instance.name\",\"aws_db_instance.db_instance\"]}},\"password\":{\"sensitive\":true,\"expression\":{\"references\":[\"aws_db_instance.db_instance.password\",\"aws_db_instance.db_instance\"]}},\"provider_verify_certificate\":{\"expression\":{\"references\":[\"var.provider_verify_certificate\"]}},\"require_ssl\":{\"expression\":{\"references\":[\"var.require_ssl\"]}},\"status\":{\"expression\":{\"references\":[\"aws_db_instance.db_instance.name\",\"aws_db_instance.db_instance\",\"aws_db_instance.db_instance.id\",\"aws_db_instance.db_instance\",\"aws_db_instance.db_instance.address\",\"aws_db_instance.db_instance\",\"var.region\",\"var.region\",\"aws_db_instance.db_instance.id\",\"aws_db_instance.db_instance\"]}},\"username\":{\"expression\":{\"references\":[\"aws_db_instance.db_instance.username\",\"aws_db_instance.db_instance\"]}}},\"resources\":[{\"address\":\"aws_db_instance.db_instance\",\"mode\":\"managed\",\"type\":\"aws_db_instance\",\"name\":\"db_instance\",\"provider_config_key\":\"aws\",\"expressions\":{\"allocated_storage\":{\"references\":[\"var.storage_gb\"]},\"allow_major_version_upgrade\":{\"references\":[\"var.allow_major_version_upgrade\"]},\"apply_immediately\":{\"constant_value\":true},\"auto_minor_version_upgrade\":{\"references\":[\"var.auto_minor_version_upgrade\"]},\"backup_retention_period\":{\"references\":[\"var.backup_retention_period\"]},\"backup_window\":{\"references\":[\"var.backup_window\",\"var.backup_window\"]},\"copy_tags_to_snapshot\":{\"references\":[\"var.copy_tags_to_snapshot\"]},\"db_name\":{\"references\":[\"var.db_name\"]},\"db_subnet_group_name\":{\"references\":[\"local.subnet_group\"]},\"delete_automated_backups\":{\"references\":[\"var.delete_automated_backups\"]},\"deletion_protection\":{\"references\":[\"var.deletion_protection\"]},\"engine\":{\"references\":[\"local.engine\"]},\"engine_version\":{\"references\":[\"var.postgres_version\"]},\"identifier\":{\"references\":[\"var.instance_name\"]},\"instance_class\":{\"references\":[\"local.instance_class\"]},\"iops\":{\"references\":[\"var.storage_type\",\"var.iops\"]},\"kms_key_id\":{\"references\":[\"var.kms_key_id\",\"var.kms_key_id\"]},\"maintenance_window\":{\"references\":[\"var.maintenance_window\",\"var.maintenance_window\"]},\"max_allocated_storage\":{\"references\":[\"local.max_allocated_storage\"]},\"monitoring_interval\":{\"references\":[\"var.monitoring_interval\"]},\"monitoring_role_arn\":{\"references\":[\"var.monitoring_role_arn\"]},\"multi_az\":{\"references\":[\"var.multi_az\"]},\"parameter_group_name\":{\"references\":[\"var.parameter_group_name\",\"aws_db_parameter_group.db_parameter_group[0].name\",\"aws_db_parameter_group.db_parameter_group[0]\",\"aws_db_parameter_group.db_parameter_group\",\"var.parameter_group_name\"]},\"password\":{\"references\":[\"random_password.password.result\",\"random_password.password\"]},\"performance_insights_enabled\":{\"references\":[\"var.performance_insights_enabled\"]},\"performance_insights_kms_key_id\":{\"references\":[\"var.performance_insights_kms_key_id\",\"var.performance_insights_kms_key_id\"]},\"publicly_accessible\":{\"references\":[\"var.publicly_accessible\"]},\"skip_final_snapshot\":{\"constant_value\":true},\"storage_encrypted\":{\"references\":[\"var.storage_encrypted\"]},\"storage_type\":{\"references\":[\"var.storage_type\"]},\"tags\":{\"references\":[\"var.labels\"]},\"username\":{\"references\":[\"random_string.username.result\",\"random_string.username\"]},\"vpc_security_group_ids\":{\"references\":[\"local.rds_vpc_security_group_ids\"]}},\"schema_version\":1},{\"address\":\"aws_db_parameter_group.db_parameter_group\",\"mode\":\"managed\",\"type\":\"aws_db_parameter_group\",\"name\":\"db_parameter_group\",\"provider_config_key\":\"aws\",\"expressions\":{\"family\":{\"references\":[\"local.engine\",\"local.major_version\"]},\"name\":{\"references\":[\"var.instance_name\"]},\"parameter\":[{\"name\":{\"constant_value\":\"rds.force_ssl\"},\"value\":{\"references\":[\"var.require_ssl\"]}}]},\"schema_version\":0,\"count_expression\":{\"references\":[\"var.parameter_group_name\"]}},{\"address\":\"aws_db_subnet_group.rds-private-subnet\",\"mode\":\"managed\",\"type\":\"aws_db_subnet_group\",\"name\":\"rds-private-subnet\",\"provider_config_key\":\"aws\",\"expressions\":{\"name\":{\"references\":[\"var.instance_name\"]},\"subnet_ids\":{\"references\":[\"data.aws_subnets.all.ids\",\"data.aws_subnets.all\"]}},\"schema_version\":0,\"count_expression\":{\"references\":[\"var.rds_subnet_group\"]}},{\"address\":\"aws_security_group.rds-sg\",\"mode\":\"managed\",\"type\":\"aws_security_group\",\"name\":\"rds-sg\",\"provider_config_key\":\"aws\",\"expressions\":{\"name\":{\"references\":[\"var.instance_name\"]},\"vpc_id\":{\"references\":[\"data.aws_vpc.vpc.id\",\"data.aws_vpc.vpc\"]}},\"schema_version\":1,\"count_expression\":{\"references\":[\"var.rds_vpc_security_group_ids\"]}},{\"address\":\"aws_security_group_rule.rds_inbound_access\",\"mode\":\"managed\",\"type\":\"aws_security_group_rule\",\"name\":\"rds_inbound_access\",\"provider_config_key\":\"aws\",\"expressions\":{\"cidr_blocks\":{\"constant_value\":[\"0.0.0.0/0\"]},\"from_port\":{\"references\":[\"local.port\"]},\"protocol\":{\"constant_value\":\"tcp\"},\"security_group_id\":{\"references\":[\"aws_security_group.rds-sg[0].id\",\"aws_security_group.rds-sg[0]\",\"aws_security_group.rds-sg\"]},\"to_port\":{\"references\":[\"local.port\"]},\"type\":{\"constant_value\":\"ingress\"}},\"schema_version\":2,\"count_expression\":{\"references\":[\"var.rds_vpc_security_group_ids\"]}},{\"address\":\"random_password.password\",\"mode\":\"managed\",\"type\":\"random_password\",\"name\":\"password\",\"provider_config_key\":\"random\",\"expressions\":{\"length\":{\"constant_value\":32},\"override_special\":{\"constant_value\":\"~_-.\"},\"special\":{\"constant_value\":false}},\"schema_version\":2},{\"address\":\"random_string.username\",\"mode\":\"managed\",\"type\":\"random_string\",\"name\":\"username\",\"provider_config_key\":\"random\",\"expressions\":{\"length\":{\"constant_value\":16},\"numeric\":{\"constant_value\":false},\"special\":{\"constant_value\":false}},\"schema_version\":2},{\"address\":\"data.aws_subnets.all\",\"mode\":\"data\",\"type\":\"aws_subnets\",\"name\":\"all\",\"provider_config_key\":\"aws\",\"expressions\":{\"filter\":[{\"name\":{\"constant_value\":\"vpc-id\"},\"values\":{\"references\":[\"data.aws_vpc.vpc.id\",\"data.aws_vpc.vpc\"]}}]},\"schema_version\":0},{\"address\":\"data.aws_vpc.vpc\",\"mode\":\"data\",\"type\":\"aws_vpc\",\"name\":\"vpc\",\"provider_config_key\":\"aws\",\"expressions\":{\"default\":{\"references\":[\"var.aws_vpc_id\"]},\"id\":{\"references\":[\"var.aws_vpc_id\",\"var.aws_vpc_id\"]}},\"schema_version\":0}],\"variables\":{\"allow_major_version_upgrade\":{},\"auto_minor_version_upgrade\":{},\"aws_access_key_id\":{},\"aws_secret_access_key\":{},\"aws_vpc_id\":{},\"backup_retention_period\":{},\"backup_window\":{},\"copy_tags_to_snapshot\":{},\"cores\":{},\"db_name\":{},\"delete_automated_backups\":{},\"deletion_protection\":{},\"instance_class\":{},\"instance_name\":{},\"iops\":{},\"kms_key_id\":{},\"labels\":{},\"maintenance_window\":{},\"monitoring_interval\":{},\"monitoring_role_arn\":{},\"multi_az\":{},\"parameter_group_name\":{},\"performance_insights_enabled\":{},\"performance_insights_kms_key_id\":{},\"postgres_version\":{},\"provider_verify_certificate\":{},\"publicly_accessible\":{},\"rds_subnet_group\":{},\"rds_vpc_security_group_ids\":{},\"region\":{},\"require_ssl\":{},\"storage_autoscale\":{},\"storage_autoscale_limit_gb\":{},\"storage_encrypted\":{},\"storage_gb\":{},\"storage_type\":{},\"subsume\":{}}}},\"relevant_attributes\":[{\"resource\":\"random_password.password\",\"attribute\":[\"result\"]},{\"resource\":\"aws_db_parameter_group.db_parameter_group[0]\",\"attribute\":[\"name\"]},{\"resource\":\"aws_db_instance.db_instance\",\"attribute\":[\"name\"]},{\"resource\":\"aws_db_instance.db_instance\",\"attribute\":[\"password\"]},{\"resource\":\"aws_db_instance.db_instance\",\"attribute\":[\"username\"]},{\"resource\":\"data.aws_vpc.vpc\",\"attribute\":[\"id\"]},{\"resource\":\"aws_security_group.rds-sg[0]\",\"attribute\":[\"id\"]},{\"resource\":\"aws_db_instance.db_instance\",\"attribute\":[\"id\"]},{\"resource\":\"aws_db_subnet_group.rds-private-subnet[0]\",\"attribute\":[\"name\"]},{\"resource\":\"random_string.username\",\"attribute\":[\"result\"]},{\"resource\":\"aws_db_instance.db_instance\",\"attribute\":[\"address\"]}]}"

	var outputJson tfjson.Plan
	fmt.Fprintf(GinkgoWriter, "load json")
	err := json.Unmarshal([]byte(json_string), &outputJson)
	fmt.Fprintf(GinkgoWriter, "%v", outputJson)
	Expect(err).To(BeNil())
}