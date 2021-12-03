/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	alibabacloudcomv1 "github.com/aoxn/ovm/pkg/apis/alibabacloud.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMasterSets implements MasterSetInterface
type FakeMasterSets struct {
	Fake *FakeOvmV1
	ns   string
}

var mastersetsResource = schema.GroupVersionResource{Group: "alibabacloud.com", Version: "v1", Resource: "mastersets"}

var mastersetsKind = schema.GroupVersionKind{Group: "alibabacloud.com", Version: "v1", Kind: "MasterSet"}

// Get takes name of the masterSet, and returns the corresponding masterSet object, and an error if there is any.
func (c *FakeMasterSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *alibabacloudcomv1.MasterSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mastersetsResource, c.ns, name), &alibabacloudcomv1.MasterSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*alibabacloudcomv1.MasterSet), err
}

// List takes label and field selectors, and returns the list of MasterSets that match those selectors.
func (c *FakeMasterSets) List(ctx context.Context, opts v1.ListOptions) (result *alibabacloudcomv1.MasterSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mastersetsResource, mastersetsKind, c.ns, opts), &alibabacloudcomv1.MasterSetList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &alibabacloudcomv1.MasterSetList{ListMeta: obj.(*alibabacloudcomv1.MasterSetList).ListMeta}
	for _, item := range obj.(*alibabacloudcomv1.MasterSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested masterSets.
func (c *FakeMasterSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mastersetsResource, c.ns, opts))

}

// Create takes the representation of a masterSet and creates it.  Returns the server's representation of the masterSet, and an error, if there is any.
func (c *FakeMasterSets) Create(ctx context.Context, masterSet *alibabacloudcomv1.MasterSet, opts v1.CreateOptions) (result *alibabacloudcomv1.MasterSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mastersetsResource, c.ns, masterSet), &alibabacloudcomv1.MasterSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*alibabacloudcomv1.MasterSet), err
}

// Update takes the representation of a masterSet and updates it. Returns the server's representation of the masterSet, and an error, if there is any.
func (c *FakeMasterSets) Update(ctx context.Context, masterSet *alibabacloudcomv1.MasterSet, opts v1.UpdateOptions) (result *alibabacloudcomv1.MasterSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mastersetsResource, c.ns, masterSet), &alibabacloudcomv1.MasterSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*alibabacloudcomv1.MasterSet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMasterSets) UpdateStatus(ctx context.Context, masterSet *alibabacloudcomv1.MasterSet, opts v1.UpdateOptions) (*alibabacloudcomv1.MasterSet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mastersetsResource, "status", c.ns, masterSet), &alibabacloudcomv1.MasterSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*alibabacloudcomv1.MasterSet), err
}

// Delete takes name of the masterSet and deletes it. Returns an error if one occurs.
func (c *FakeMasterSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mastersetsResource, c.ns, name), &alibabacloudcomv1.MasterSet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMasterSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mastersetsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &alibabacloudcomv1.MasterSetList{})
	return err
}

// Patch applies the patch and returns the patched masterSet.
func (c *FakeMasterSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *alibabacloudcomv1.MasterSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mastersetsResource, c.ns, name, pt, data, subresources...), &alibabacloudcomv1.MasterSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*alibabacloudcomv1.MasterSet), err
}