package appendstorage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/minio-io/minio/pkgs/storage"
	. "gopkg.in/check.v1"
)

type AppendStorageSuite struct{}

var _ = Suite(&AppendStorageSuite{})

func Test(t *testing.T) { TestingT(t) }

func makeTempTestDir() (string, error) {
	return ioutil.TempDir("/tmp", "minio-test-")
}

func (s *AppendStorageSuite) TestAppendStoragePutAtRootPath(c *C) {
	rootDir, err := makeTempTestDir()
	c.Assert(err, IsNil)
	defer os.RemoveAll(rootDir)

	var objectStorage storage.ObjectStorage
	objectStorage, err = NewStorage(rootDir, 0)
	c.Assert(err, IsNil)

	err = objectStorage.Put("path1", []byte("object1"))
	c.Assert(err, IsNil)

	// assert object1 was created in correct path
	object1, err := objectStorage.Get("path1")
	c.Assert(err, IsNil)
	c.Assert(string(object1), Equals, "object1")

	err = objectStorage.Put("path2", []byte("object2"))
	c.Assert(err, IsNil)

	// assert object1 was created in correct path
	object2, err := objectStorage.Get("path2")
	c.Assert(err, IsNil)
	c.Assert(string(object2), Equals, "object2")

	object1, err = objectStorage.Get("path1")
	c.Assert(err, IsNil)
	c.Assert(string(object1), Equals, "object1")
}

func (s *AppendStorageSuite) TestAppendStoragePutDirPath(c *C) {
	rootDir, err := makeTempTestDir()
	c.Assert(err, IsNil)
	defer os.RemoveAll(rootDir)

	var objectStorage storage.ObjectStorage
	objectStorage, err = NewStorage(rootDir, 0)
	c.Assert(err, IsNil)

	// add object 1
	objectStorage.Put("path1/path2/path3", []byte("object"))

	// assert object1 was created in correct path
	object1, err := objectStorage.Get("path1/path2/path3")
	c.Assert(err, IsNil)
	c.Assert(string(object1), Equals, "object")

	// add object 2
	objectStorage.Put("path1/path1/path1", []byte("object2"))

	// assert object1 was created in correct path
	object2, err := objectStorage.Get("path1/path1/path1")
	c.Assert(err, IsNil)
	c.Assert(string(object2), Equals, "object2")
}

func (s *AppendStorageSuite) TestSerialization(c *C) {
	rootDir, err := makeTempTestDir()
	c.Assert(err, IsNil)
	defer os.RemoveAll(rootDir)

	objectStorage, err := NewStorage(rootDir, 0)
	c.Assert(err, IsNil)

	err = objectStorage.Put("path1", []byte("object1"))
	c.Assert(err, IsNil)
	err = objectStorage.Put("path2", []byte("object2"))
	c.Assert(err, IsNil)
	err = objectStorage.Put("path3/obj3", []byte("object3"))
	c.Assert(err, IsNil)

	es := objectStorage.(*appendStorage)
	es.file.Close()

	objectStorage2, err := NewStorage(rootDir, 0)
	c.Assert(err, IsNil)

	object1, err := objectStorage2.Get("path1")
	c.Assert(err, IsNil)
	c.Assert(string(object1), Equals, "object1")

	object2, err := objectStorage2.Get("path2")
	c.Assert(err, IsNil)
	c.Assert(string(object2), Equals, "object2")

	object3, err := objectStorage2.Get("path3/obj3")
	c.Assert(err, IsNil)
	c.Assert(string(object3), Equals, "object3")
}

func (s *AppendStorageSuite) TestSlice(c *C) {
}