package data_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielboakye/go-echo-app/data"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

func newTestRepo() (
	sqlmock.Sqlmock, data.IRepository,
) {
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	Expect(err).Should(BeNil())

	testRepo := data.NewRepository(db)

	return mockDB, testRepo
}
