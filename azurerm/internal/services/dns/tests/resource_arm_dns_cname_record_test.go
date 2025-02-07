package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDnsCNameRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	config := testAccAzureRMDnsCNameRecord_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCNameRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCNameRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsCNameRecord_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_cname_record"),
			},
		},
	})
}

func TestAccAzureRMDnsCNameRecord_subdomain(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	config := testAccAzureRMDnsCNameRecord_subdomain(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record", "test.contoso.com"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCNameRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	preConfig := testAccAzureRMDnsCNameRecord_basic(data)
	postConfig := testAccAzureRMDnsCNameRecord_updateRecords(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDnsCNameRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	preConfig := testAccAzureRMDnsCNameRecord_withTags(data)
	postConfig := testAccAzureRMDnsCNameRecord_withTagsUpdate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCNameRecord_withAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	targetResourceName := "azurerm_dns_cname_record.target"
	targetResourceName2 := "azurerm_dns_cname_record.target2"
	preConfig := testAccAzureRMDnsCNameRecord_withAlias(data)
	postConfig := testAccAzureRMDnsCNameRecord_withAliasUpdate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName2, "id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCNameRecord_RecordToAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	targetResourceName := "azurerm_dns_cname_record.target2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCNameRecord_AliasToRecordUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMDnsCNameRecord_AliasToRecord(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "record"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCNameRecord_AliasToRecord(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record", "test")
	targetResourceName := "azurerm_dns_cname_record.target2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCNameRecord_AliasToRecord(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: testAccAzureRMDnsCNameRecord_AliasToRecordUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "target_resource_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsCNameRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		cnameName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS CNAME record: %s", cnameName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, cnameName, dns.CNAME)
		if err != nil {
			return fmt.Errorf("Bad: Get CNAME RecordSet: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS CNAME record %s (resource group: %s) does not exist", cnameName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsCNameRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_cname_record" {
			continue
		}

		cnameName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, cnameName, dns.CNAME)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS CNAME record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsCNameRecord_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsCNameRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_cname_record" "import" {
  name                = "${azurerm_dns_cname_record.test.name}"
  resource_group_name = "${azurerm_dns_cname_record.test.resource_group_name}"
  zone_name           = "${azurerm_dns_cname_record.test.zone_name}"
  ttl                 = 300
  record              = "contoso.com"
}
`, template)
}

func testAccAzureRMDnsCNameRecord_subdomain(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "test.contoso.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.co.uk"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_withAlias(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "target" {
	name                = "mycnametarget%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	zone_name           = "${azurerm_dns_zone.test.name}"
	ttl                 = 300
	record              = "contoso.com"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "mycnamerecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  target_resource_id   = "${azurerm_dns_cname_record.target.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_withAliasUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "target2" {
	name                = "mycnametarget%d2"
	resource_group_name = "${azurerm_resource_group.test.name}"
	zone_name           = "${azurerm_dns_zone.test.name}"
	ttl                 = 300
	record              = "contoso.co.uk"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "mycnamerecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  target_resource_id  = "${azurerm_dns_cname_record.target2.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_AliasToRecord(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "target2" {
	name                = "mycnametarget%d2"
	resource_group_name = "${azurerm_resource_group.test.name}"
	zone_name           = "${azurerm_dns_zone.test.name}"
	ttl                 = 300
	record              = "contoso.co.uk"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "mycnamerecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  target_resource_id  = "${azurerm_dns_cname_record.target2.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCNameRecord_AliasToRecordUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  record              = "1.2.3.4"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
