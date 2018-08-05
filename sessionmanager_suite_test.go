package sessionmanager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSessionmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sessionmanager Suite")
}
