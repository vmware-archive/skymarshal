package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/auth/authfakes"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckBuildWriteAccessHandler", func() {
	var (
		response       *http.Response
		server         *httptest.Server
		delegate       *buildDelegateHandler
		buildsDB       *authfakes.FakeBuildsDB
		handlerFactory auth.CheckBuildWriteAccessHandlerFactory
		handler        http.Handler

		authValidator     *authfakes.FakeValidator
		userContextReader *authfakes.FakeUserContextReader

		build    *dbfakes.FakeBuild
		pipeline db.SavedPipeline
	)

	BeforeEach(func() {
		buildsDB = new(authfakes.FakeBuildsDB)
		handlerFactory = auth.NewCheckBuildWriteAccessHandlerFactory(buildsDB)

		authValidator = new(authfakes.FakeValidator)
		userContextReader = new(authfakes.FakeUserContextReader)

		delegate = &buildDelegateHandler{}

		build = new(dbfakes.FakeBuild)
		pipeline = db.SavedPipeline{}
		build.GetPipelineReturns(pipeline, nil)
		build.TeamNameReturns("some-team")
		build.JobNameReturns("some-job")

		checkBuildWriteAccessHandler := handlerFactory.HandlerFor(delegate, auth.UnauthorizedRejector{})
		handler = auth.WrapHandler(checkBuildWriteAccessHandler, authValidator, userContextReader)
	})

	JustBeforeEach(func() {
		server = httptest.NewServer(handler)

		request, err := http.NewRequest("POST", server.URL+"?:team_name=some-team&:build_id=55", nil)
		Expect(err).NotTo(HaveOccurred())

		response, err = new(http.Client).Do(request)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = AfterEach(func() {
		server.Close()
	})

	Context("when authenticated and accessing same team's build", func() {
		BeforeEach(func() {
			authValidator.IsAuthenticatedReturns(true)
			userContextReader.GetTeamReturns("some-team", 42, true, true)
		})

		Context("when build exists", func() {
			BeforeEach(func() {
				buildsDB.GetBuildByIDReturns(build, true, nil)
			})

			It("returns 200 ok", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("calls delegate with the build context", func() {
				Expect(delegate.IsCalled).To(BeTrue())
				Expect(delegate.ContextBuild).To(BeIdenticalTo(build))
			})
		})

		Context("when build is not found", func() {
			BeforeEach(func() {
				buildsDB.GetBuildByIDReturns(nil, false, nil)
			})

			It("returns 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when getting build fails", func() {
			BeforeEach(func() {
				buildsDB.GetBuildByIDReturns(nil, false, errors.New("disaster"))
			})

			It("returns 404", func() {
				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("when authenticated but accessing different team's build", func() {
		BeforeEach(func() {
			authValidator.IsAuthenticatedReturns(true)
			userContextReader.GetTeamReturns("other-team-name", 0, false, true)
			buildsDB.GetBuildByIDReturns(build, true, nil)
		})

		It("returns 403", func() {
			Expect(response.StatusCode).To(Equal(http.StatusForbidden))
		})
	})

	Context("when not authenticated", func() {
		BeforeEach(func() {
			authValidator.IsAuthenticatedReturns(false)
			userContextReader.GetTeamReturns("", 0, false, false)
		})

		It("returns 401", func() {
			Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})
