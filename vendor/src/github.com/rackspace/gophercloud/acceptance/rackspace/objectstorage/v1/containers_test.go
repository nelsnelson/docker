// +build acceptance rackspace objectstorage v1

package v1

import (
	"testing"

	osContainers "github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/pagination"
	raxContainers "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestContainers(t *testing.T) {
	c, err := createClient(t, false)
	th.AssertNoErr(t, err)

	t.Logf("Containers Info available to the currently issued token:")
	count := 0
	err = raxContainers.List(c, &osContainers.ListOpts{Full: true}).EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page %02d ---", count)

		containers, err := raxContainers.ExtractInfo(page)
		th.AssertNoErr(t, err)

		for i, container := range containers {
			t.Logf("[%02d]      name=[%s]", i, container.Name)
			t.Logf("            count=[%d]", container.Count)
			t.Logf("            bytes=[%d]", container.Bytes)
		}

		count++
		return true, nil
	})
	th.AssertNoErr(t, err)
	if count == 0 {
		t.Errorf("No containers listed for your current token.")
	}

	t.Logf("Container Names available to the currently issued token:")
	count = 0
	err = raxContainers.List(c, &osContainers.ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page %02d ---", count)

		names, err := raxContainers.ExtractNames(page)
		th.AssertNoErr(t, err)

		for i, name := range names {
			t.Logf("[%02d] %s", i, name)
		}

		count++
		return true, nil
	})
	th.AssertNoErr(t, err)
	if count == 0 {
		t.Errorf("No containers listed for your current token.")
	}

	headers, err := raxContainers.Create(c, "gophercloud-test", nil).ExtractHeaders()
	th.AssertNoErr(t, err)
	defer func() {
		_, err := raxContainers.Delete(c, "gophercloud-test").ExtractHeaders()
		th.AssertNoErr(t, err)
	}()

	headers, err = raxContainers.Update(c, "gophercloud-test", raxContainers.UpdateOpts{Metadata: map[string]string{"white": "mountains"}}).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Update Account request: %+v\n", headers)
	defer func() {
		_, err := raxContainers.Update(c, "gophercloud-test", raxContainers.UpdateOpts{Metadata: map[string]string{"white": ""}}).ExtractHeaders()
		th.AssertNoErr(t, err)
		metadata, err := raxContainers.Get(c, "gophercloud-test").ExtractMetadata()
		th.AssertNoErr(t, err)
		t.Logf("Metadata from Get Account request (after update reverted): %+v\n", metadata)
		th.CheckEquals(t, metadata["White"], "")
	}()

	getResult := raxContainers.Get(c, "gophercloud-test")
	headers, err = getResult.ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Get Account request (after update): %+v\n", headers)
	metadata, err := getResult.ExtractMetadata()
	th.AssertNoErr(t, err)
	t.Logf("Metadata from Get Account request (after update): %+v\n", metadata)
	th.CheckEquals(t, metadata["White"], "mountains")
}