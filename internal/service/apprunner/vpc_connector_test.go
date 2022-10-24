package apprunner_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/apprunner"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapprunner "github.com/hashicorp/terraform-provider-aws/internal/service/apprunner"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccAppRunnerVPCConnector_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_vpc_connector.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckVPCConnector(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPCConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCConnectorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectorExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`vpcconnector/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "vpc_connector_name", rName),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAppRunnerVPCConnector_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_vpc_connector.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckVPCConnector(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPCConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCConnectorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectorExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfapprunner.ResourceVPCConnector(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAppRunnerVPCConnector_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_vpc_connector.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckVPCConnector(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPCConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCConnectorConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCConnectorConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCConnectorConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckVPCConnectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_apprunner_vpc_connector" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppRunnerConn

		_, err := tfapprunner.FindVPCConnectorByARN(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("App Runner VPC Connector %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckVPCConnectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Runner VPC Connector ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppRunnerConn

		_, err := tfapprunner.FindVPCConnectorByARN(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccVPCConnectorConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.1.1.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCConnectorConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccVPCConnectorConfig_base(rName), fmt.Sprintf(`
resource "aws_apprunner_vpc_connector" "test" {
  vpc_connector_name = %[1]q
  subnets            = [aws_subnet.test.id]
  security_groups    = [aws_security_group.test.id]
}
`, rName))
}

func testAccVPCConnectorConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccVPCConnectorConfig_base(rName), fmt.Sprintf(`
resource "aws_apprunner_vpc_connector" "test" {
  vpc_connector_name = %[1]q
  subnets            = [aws_subnet.test.id]
  security_groups    = [aws_security_group.test.id]

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccVPCConnectorConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccVPCConnectorConfig_base(rName), fmt.Sprintf(`
resource "aws_apprunner_vpc_connector" "test" {
  vpc_connector_name = %[1]q
  subnets            = [aws_subnet.test.id]
  security_groups    = [aws_security_group.test.id]

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccPreCheckVPCConnector(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AppRunnerConn
	ctx := context.Background()

	input := &apprunner.ListVpcConnectorsInput{}

	_, err := conn.ListVpcConnectorsWithContext(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
