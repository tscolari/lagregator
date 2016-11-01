package lagregator_test

import (
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/lostmars/lagregator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/glager"
)

var _ = Describe("Relogger", func() {
	var (
		destLogger lager.Logger
		srcLogger  lager.Logger
	)

	BeforeEach(func() {
		destLogger = glager.NewLogger("destination")
		srcLogger = glager.NewLogger("source")
	})

	It("returns an io.Writer that relogs to destination", func() {
		relogger := lagregator.NewRelogger(destLogger)
		srcSink := lager.NewWriterSink(relogger, lager.DEBUG)
		srcLogger.RegisterSink(srcSink)

		srcLogger.Debug("first-debug", lager.Data{"attr1": "value1"})
		srcLogger.Info("first-info", lager.Data{"attr1": "value1"})
		srcLogger.Error("first-error", errors.New("failed!"), lager.Data{"attr1": "value1"})
		srcLogger.Debug("second-debug", lager.Data{"attr2": "value2"})

		Expect(destLogger).To(glager.HaveLogged(
			glager.Debug(
				glager.Source("destination"),
				glager.Message("destination.source.first-debug"),
				glager.Data("attr1", "value1"),
			),
			glager.Info(
				glager.Source("destination"),
				glager.Message("destination.source.first-info"),
				glager.Data("attr1", "value1"),
			),
			glager.Error(
				errors.New("failed!"),
				glager.Source("destination"),
				glager.Message("destination.source.first-error"),
				glager.Data("attr1", "value1"),
			),
			glager.Debug(
				glager.Source("destination"),
				glager.Message("destination.source.second-debug"),
				glager.Data("attr2", "value2"),
			),
		))
	})
})
