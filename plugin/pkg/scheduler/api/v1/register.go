/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package v1

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

// SchemeGroupVersion is group version used to register these objects
// TODO this should be in the "scheduler" group
var SchemeGroupVersion = unversioned.GroupVersion{Group: "", Version: "v1"}

// Codec encodes internal objects to the v1 scheme
var Codec = runtime.CodecFor(api.Scheme, SchemeGroupVersion.String())

func init() {
	api.Scheme.AddKnownTypes(SchemeGroupVersion,
		&Policy{},
	)
}

func (*Policy) IsAnAPIObject() {}
