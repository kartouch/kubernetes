/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package fake

import (
	api "k8s.io/kubernetes/pkg/api"
	core "k8s.io/kubernetes/pkg/client/testing/core"
	labels "k8s.io/kubernetes/pkg/labels"
	watch "k8s.io/kubernetes/pkg/watch"
)

// FakeServices implements ServiceInterface
type FakeServices struct {
	Fake *FakeLegacy
	ns   string
}

func (c *FakeServices) Create(service *api.Service) (result *api.Service, err error) {
	obj, err := c.Fake.
		Invokes(core.NewCreateAction("services", c.ns, service), &api.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.Service), err
}

func (c *FakeServices) Update(service *api.Service) (result *api.Service, err error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateAction("services", c.ns, service), &api.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.Service), err
}

func (c *FakeServices) UpdateStatus(service *api.Service) (*api.Service, error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateSubresourceAction("services", "status", c.ns, service), &api.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.Service), err
}

func (c *FakeServices) Delete(name string, options *api.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(core.NewDeleteAction("services", c.ns, name), &api.Service{})

	return err
}

func (c *FakeServices) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	action := core.NewDeleteCollectionAction("events", c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &api.ServiceList{})
	return err
}

func (c *FakeServices) Get(name string) (result *api.Service, err error) {
	obj, err := c.Fake.
		Invokes(core.NewGetAction("services", c.ns, name), &api.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.Service), err
}

func (c *FakeServices) List(opts api.ListOptions) (result *api.ServiceList, err error) {
	obj, err := c.Fake.
		Invokes(core.NewListAction("services", c.ns, opts), &api.ServiceList{})

	if obj == nil {
		return nil, err
	}

	label := opts.LabelSelector
	if label == nil {
		label = labels.Everything()
	}
	list := &api.ServiceList{}
	for _, item := range obj.(*api.ServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested services.
func (c *FakeServices) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(core.NewWatchAction("services", c.ns, opts))

}
