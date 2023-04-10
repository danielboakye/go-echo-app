package data_test

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielboakye/go-echo-app/data"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get user", func() {

	var (
		u   *data.User
		err error
		uid string
	)

	When("there is a match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"user_id", "email", "first_name",
							"last_name", "password", "user_active",
							"created_at", "updated_at",
						},
					).
						AddRow(
							uid, "example@mail.com", "Clark",
							"Kent", "$2a$12$4P.DPHoR0ULhMVCRSa8qg.HVagvaPoYG3Di9i253G9ILIli3sTGwy", 1,
							time.Now(), time.Now(),
						),
				)

			u, err = testRepo.GetOne(uid)

		})

		It("should return data", func() {
			Expect(u).ToNot(BeNil())
		})

		It("should populate the fields correctly", func() {
			Expect(u.FirstName).To(Equal("Clark"))
			Expect(u.ID).To(Equal(uid))
		})

		It("should not error", func() {
			Expect(err).To(BeNil())
		})
	})

	When("there's no match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnError(sql.ErrNoRows)

			u, err = testRepo.GetOne(uid)

		})

		It("should not return data", func() {
			Expect(u).To(BeNil())
		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})
})

var _ = Describe("Get All Users", func() {

	var (
		users []*data.User
		err   error
	)

	When("there is a match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, user_active, created_at, updated_at
						FROM users
						ORDER BY last_name ASC
					`).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"user_id", "email", "first_name",
							"last_name", "user_active",
							"created_at", "updated_at",
						},
					).
						AddRow(
							"2899bacc-7107-4cd4-9364-6a6fc4fc2fd3", "example@mail.com", "Clark",
							"Kent", 1,
							time.Now(), time.Now(),
						).
						AddRow(
							"ae17b2e2-6b87-4c5b-9c94-3623dacf113b", "example1@mail.com", "Lois",
							"Lane", 0,
							time.Now(), time.Now(),
						),
				)

			users, err = testRepo.GetAll()

		})

		It("should return data", func() {
			Expect(users).ToNot(BeNil())
			Expect(users).To(HaveLen(2))
		})

		It("should populate the fields correctly", func() {
			Expect(users[0].ID).To(Equal("2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"))
			Expect(users[1].Active).To(Equal(0))
			Expect(users[0].Email).ToNot(Equal(users[1].Email))
		})

		It("should not error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("there's no match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, user_active, created_at, updated_at
						FROM users
						ORDER BY last_name ASC
					`).
				WillReturnError(sql.ErrNoRows)

			users, err = testRepo.GetAll()

		})

		It("should not return data", func() {
			Expect(users).To(BeNil())
		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})
})

var _ = Describe("Get user by Email", func() {

	var (
		u     *data.User
		err   error
		email string
	)

	When("there is a match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			email = "emaple@gmail.com"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE email = $1
					`).
				WithArgs(email).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"user_id", "email", "first_name",
							"last_name", "password", "user_active",
							"created_at", "updated_at",
						},
					).
						AddRow(
							"2899bacc-7107-4cd4-9364-6a6fc4fc2fd3", email, "Clark",
							"Kent", "$2a$12$4P.DPHoR0ULhMVCRSa8qg.HVagvaPoYG3Di9i253G9ILIli3sTGwy", 1,
							time.Now(), time.Now(),
						),
				)

			u, err = testRepo.GetByEmail(email)

		})

		It("should return data", func() {
			Expect(u).ToNot(BeNil())
		})

		It("should populate the fields correctly", func() {
			Expect(u.FirstName).To(Equal("Clark"))
			Expect(u.Email).To(Equal(email))
		})

		It("should not error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("there's no match", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			email = "emaple-notexist@gmail.com"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE email = $1
					`).
				WithArgs(email).
				WillReturnError(sql.ErrNoRows)

			u, err = testRepo.GetByEmail(email)

		})

		It("should not return data", func() {
			Expect(u).To(BeNil())
		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})
})

var _ = Describe("Update user", func() {

	var (
		err error
		u   data.User
	)

	When("successfully updated", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			u = data.User{
				ID:        "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3",
				Email:     "example@gmail.com",
				FirstName: "Clark",
				LastName:  "Kent",
				Active:    1,
				UpdatedAt: time.Now(),
			}

			mockDB.ExpectExec(`
						UPDATE users
						SET 
							email = $1, first_name = $2,
							last_name = $3, updated_at = $4, 
							user_active = $5
						WHERE user_id = $6
					`).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = testRepo.Update(u)

		})

		It("should not error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("it fails", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			u = data.User{
				ID:        "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3",
				Email:     "example@gmail.com",
				FirstName: "Clark",
				LastName:  "Kent",
				Active:    1,
				UpdatedAt: time.Now(),
			}

			mockDB.ExpectExec(`
						UPDATE users
						SET
							email = $1, first_name = $2,
							last_name = $3, updated_at = $4,
							user_active = $5
						WHERE user_id = $6
					`).
				WillReturnError(sql.ErrConnDone)

			err = testRepo.Update(u)

		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrConnDone))
		})
	})
})

var _ = Describe("Delete user by id", func() {

	var (
		err error
		uid string
	)

	When("successfully deleted", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"

			mockDB.ExpectExec(`
						DELETE FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = testRepo.DeleteByID(uid)

		})

		It("should not error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("it fails", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"

			mockDB.ExpectExec(`
						DELETE FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnError(sql.ErrConnDone)

			err = testRepo.DeleteByID(uid)

		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrConnDone))
		})
	})
})

var _ = Describe("Insert user into db", func() {

	var (
		err error
		u   data.User
		uid string
		rid string
	)

	When("successful", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"
			u = data.User{
				Email:     "example@gmail.com",
				FirstName: "Clark",
				LastName:  "Kent",
				Password:  "$2a$12$4P.DPHoR0ULhMVCRSa8qg.HVagvaPoYG3Di9i253G9ILIli3sTGwy",
				Active:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			mockDB.ExpectQuery(`
					INSERT INTO users (email,first_name,last_name,password,user_active,created_at,updated_at) 
					VALUES ($1,$2,$3,$4,$5,$6,$7) 
					RETURNING user_id`,
			).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uid))

			rid, err = testRepo.Insert(u)

		})

		It("should not error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return user id", func() {
			Expect(rid).To(Equal(uid))
		})

	})

	When("it fails", func() {
		BeforeEach(func() {
			mockDB, testRepo := newTestRepo()
			uid = "2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"
			u = data.User{
				Email:     "example@gmail.com",
				FirstName: "Clark",
				LastName:  "Kent",
				Password:  "$2a$12$4P.DPHoR0ULhMVCRSa8qg.HVagvaPoYG3Di9i253G9ILIli3sTGwy",
				Active:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			mockDB.ExpectQuery(`
					INSERT INTO users (email,first_name,last_name,password,user_active,created_at,updated_at) 
					VALUES ($1,$2,$3,$4,$5,$6,$7) 
					RETURNING user_id`,
			).
				WillReturnError(sql.ErrConnDone)

			rid, err = testRepo.Insert(u)

		})

		It("should return error", func() {
			Expect(err).To(MatchError(sql.ErrConnDone))
		})

		It("should return user id", func() {
			Expect(rid).ShouldNot(Equal(uid))
		})
	})
})
