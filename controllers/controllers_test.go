package controllers_test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielboakye/go-echo-app/data"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get all users", func() {

	var (
		err  error
		body []byte
		resp *http.Response
		u    []data.User
	)

	Context("successful request", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()

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

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/users", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			err = json.Unmarshal(body, &u)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusOK).To(Equal(resp.StatusCode))
		})

		It("should populate the fields correctly", func() {
			Expect(u[0].ID).To(Equal("2899bacc-7107-4cd4-9364-6a6fc4fc2fd3"))
			Expect(u[1].Active).To(Equal(0))
			Expect(u[0].Email).ToNot(Equal(u[1].Email))
		})

	})

	When("request fails", func() {
		BeforeEach(func() {
			app, _ := newTestApp()

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/users", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})

	})
})

var _ = Describe("get User", func() {

	var (
		err  error
		body []byte
		resp *http.Response
		u    data.User
		uid  string
	)

	Context("successful request", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "ae17b2e2-6b87-4c5b-9c94-3623dacf113b"

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

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/users/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			err = json.Unmarshal(body, &u)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusOK).To(Equal(resp.StatusCode))
		})

		It("should populate the fields correctly", func() {
			Expect(u.ID).To(Equal(uid))
			Expect(u.Active).To(Equal(1))
		})

	})

	When("request fails", func() {
		BeforeEach(func() {
			app, _ := newTestApp()
			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/users/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})

	})
})

var _ = Describe("save User", func() {

	var (
		err   error
		body  []byte
		resp  *http.Response
		u     data.User
		uid   string
		email string
	)

	Context("successful request", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "ae17b2e2-6b87-4c5b-9c94-3623dacf113b"
			email = "example@gmail.com"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE email = $1
					`).
				WithArgs(email).
				WillReturnError(sql.ErrNoRows)

			mockDB.ExpectQuery(`
					INSERT INTO users (email,first_name,last_name,password,user_active,created_at,updated_at) 
					VALUES ($1,$2,$3,$4,$5,$6,$7) 
					RETURNING user_id`,
			).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uid))

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users", strings.NewReader(`
				{
					"email": "example@gmail.com",
					"first_name": "test",
					"last_name": "test",
					"password": "password",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			err = json.Unmarshal(body, &u)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusCreated).To(Equal(resp.StatusCode))
		})

		It("should populate the fields correctly", func() {
			Expect(u.ID).To(Equal(uid))
			Expect(u.Email).To(Equal(email))
			Expect(u.LastName).To(Equal("test"))
		})

	})

	When("email exists", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			email = "example-exists@gmail.com"

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

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users", strings.NewReader(`
				{
					"email": "example-exists@gmail.com",
					"first_name": "test",
					"last_name": "test",
					"password": "password",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})

		It("should send the correct response", func() {
			replacer := strings.NewReplacer("\r", "", "\n", "")
			Expect(replacer.Replace(string(body))).To(Equal(`{"error":"user exists"}`))
		})
	})

	Context("request fails - database insert error", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE email = $1
					`).
				WithArgs(email).
				WillReturnError(sql.ErrNoRows)

			mockDB.ExpectQuery(`
					INSERT INTO users (email,first_name,last_name,password,user_active,created_at,updated_at) 
					VALUES ($1,$2,$3,$4,$5,$6,$7) 
					RETURNING user_id`,
			).
				WillReturnError(sql.ErrConnDone)

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users", strings.NewReader(`
				{
					"email": "example-exists@gmail.com",
					"first_name": "test",
					"last_name": "test",
					"password": "password",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})

		It("should send the correct response", func() {
			replacer := strings.NewReplacer("\r", "", "\n", "")
			Expect(replacer.Replace(string(body))).To(Equal(`{"error":"submission failed"}`))
		})
	})

	When("invalid request", func() {
		BeforeEach(func() {
			app, _ := newTestApp()

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users", strings.NewReader(``))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})
	})

	When("db connection error", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			email = "example@gmail.com"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE email = $1
					`).
				WithArgs(email).
				WillReturnError(sql.ErrConnDone)

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users", strings.NewReader(`
				{
					"email": "example@gmail.com",
					"first_name": "test",
					"last_name": "test",
					"password": "password",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

			body, err = io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})

		It("should send the correct response", func() {
			replacer := strings.NewReplacer("\r", "", "\n", "")
			Expect(replacer.Replace(string(body))).To(Equal(`{"error":"processing error"}`))
		})
	})
})

var _ = Describe("update User", func() {

	var (
		resp *http.Response
		uid  string
	)

	Context("successful request", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

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

			mockDB.ExpectExec(`
						UPDATE users
						SET
							email = $1, first_name = $2,
							last_name = $3, updated_at = $4,
							user_active = $5
						WHERE user_id = $6
					`).
				WillReturnResult(sqlmock.NewResult(1, 1))

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users/61296308-2148-463d-b888-1010b3d9643b", strings.NewReader(`
				{
					"email": "example@gmail.com",
					"first_name": "Lois",
					"last_name": "Lane",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusAccepted).To(Equal(resp.StatusCode))
		})
	})

	When("update fails", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

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

			mockDB.ExpectExec(`
						UPDATE users
						SET
							email = $1, first_name = $2,
							last_name = $3, updated_at = $4,
							user_active = $5
						WHERE user_id = $6
					`).
				WillReturnError(sql.ErrConnDone)

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users/61296308-2148-463d-b888-1010b3d9643b", strings.NewReader(`
				{
					"email": "example@gmail.com",
					"first_name": "Lois",
					"last_name": "Lane",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})
	})

	When("user does not exist", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

			mockDB.ExpectQuery(`
						SELECT 
							user_id, email, first_name, last_name, password, user_active, created_at, updated_at
						FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnError(sql.ErrNoRows)

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users/61296308-2148-463d-b888-1010b3d9643b", strings.NewReader(`
				{
					"email": "example@gmail.com",
					"first_name": "Lois",
					"last_name": "Lane",
					"active": 1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})
	})

	When("invalid request", func() {
		BeforeEach(func() {
			app, _ := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/users/61296308-2148-463d-b888-1010b3d9643b", strings.NewReader(``))
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})
	})
})

var _ = Describe("delete User", func() {

	var (
		resp *http.Response
		uid  string
	)

	Context("successful request", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

			mockDB.ExpectExec(`
						DELETE FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnResult(sqlmock.NewResult(1, 1))

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http:/users/61296308-2148-463d-b888-1010b3d9643b", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusAccepted).To(Equal(resp.StatusCode))
		})
	})

	When("update fails", func() {
		BeforeEach(func() {
			app, mockDB := newTestApp()
			uid = "61296308-2148-463d-b888-1010b3d9643b"

			mockDB.ExpectExec(`
						DELETE FROM users
						WHERE user_id = $1
					`).
				WithArgs(uid).
				WillReturnError(sql.ErrConnDone)

			e := app.NewServer()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http:/users/61296308-2148-463d-b888-1010b3d9643b", nil)
			r.Header.Set("Content-Type", "application/json")
			e.ServeHTTP(w, r)

			resp = w.Result()

		})

		It("should set correct status code", func() {
			Expect(http.StatusBadRequest).To(Equal(resp.StatusCode))
		})
	})
})
