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

package volume

import (
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"os"
	"path"
)

// Volume represents a directory used by pods or hosts on a node.
// All method implementations of methods in the volume interface must be idempotent.
type Volume interface {
	// GetPath returns the directory path the volume is mounted to.
	GetPath() string

	// MetricsProvider embeds methods for exposing metrics (e.g. used,available space).
	MetricsProvider
}

// MetricsProvider exposes metrics (e.g. used,available space) related to a Volume.
type MetricsProvider interface {
	// GetMetrics returns the Metrics for the Volume.  Maybe expensive for some implementations.
	GetMetrics() (*Metrics, error)
}

// Metrics represents the used and available bytes of the Volume.
type Metrics struct {
	// Used represents the total bytes used by the Volume.
	// Note: For block devices this maybe more than the total size of the files.
	Used *resource.Quantity

	// Capacity represents the total capacity (bytes) of the volume's underlying storage.
	// For Volumes that share a filesystem with the host (e.g. emptydir, hostpath) this is the size
	// of the underlying storage, and will not equal Used + Available as the fs is shared.
	Capacity *resource.Quantity

	// Available represents the storage space available (bytes) for the Volume.
	// For Volumes that share a filesystem with the host (e.g. emptydir, hostpath), this is the available
	// space on the underlying storage, and is shared with host processes and other Volumes.
	Available *resource.Quantity
}

// Attributes represents the attributes of this builder.
type Attributes struct {
	ReadOnly                    bool
	Managed                     bool
	SupportsOwnershipManagement bool
	SupportsSELinux             bool
}

// Builder interface provides methods to set up/mount the volume.
type Builder interface {
	// Uses Interface to provide the path for Docker binds.
	Volume
	// SetUp prepares and mounts/unpacks the volume to a self-determined
	// directory path.  This may be called more than once, so
	// implementations must be idempotent.
	SetUp() error
	// SetUpAt prepares and mounts/unpacks the volume to the specified
	// directory path, which may or may not exist yet.  This may be called
	// more than once, so implementations must be idempotent.
	SetUpAt(dir string) error
	// GetAttributes returns the attributes of the builder.
	GetAttributes() Attributes
}

// Cleaner interface provides methods to cleanup/unmount the volumes.
type Cleaner interface {
	Volume
	// TearDown unmounts the volume from a self-determined directory and
	// removes traces of the SetUp procedure.
	TearDown() error
	// TearDown unmounts the volume from the specified directory and
	// removes traces of the SetUp procedure.
	TearDownAt(dir string) error
}

// Recycler provides methods to reclaim the volume resource.
type Recycler interface {
	Volume
	// Recycle reclaims the resource.  Calls to this method should block until the recycling task is complete.
	// Any error returned indicates the volume has failed to be reclaimed.  A nil return indicates success.
	Recycle() error
}

// Create adds a new resource in the storage provider and creates a PersistentVolume for the new resource.
// Calls to Create should block until complete.
type Creater interface {
	Create() (*api.PersistentVolume, error)
}

// Delete removes the resource from the underlying storage provider.  Calls to this method should block until
// the deletion is complete. Any error returned indicates the volume has failed to be reclaimed.
// A nil return indicates success.
type Deleter interface {
	Volume
	Delete() error
}

func RenameDirectory(oldPath, newName string) (string, error) {
	newPath, err := ioutil.TempDir(path.Dir(oldPath), newName)
	if err != nil {
		return "", err
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}
