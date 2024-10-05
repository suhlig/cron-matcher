package cron_matcher_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	cron "github.com/suhlig/cron-matcher"
)

var _ = Describe("Cron", func() {
	var (
		expression string
		date       time.Time
		matches    bool
		err        error
	)

	BeforeEach(func() {
		date = time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	})

	JustBeforeEach(func() {
		matches, err = cron.Matches(expression, date)
	})

	Context("exact match", func() {
		BeforeEach(func() {
			expression = "30 10 15 5 1"
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})

	Context("wildcard", func() {
		BeforeEach(func() {
			expression = "* * * * *"
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})

	Context("range match", func() {
		BeforeEach(func() {
			expression = "0-30 10-11 15-16 5-6 1-2"
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})

	Context("step match", func() {
		BeforeEach(func() {
			expression = "*/15 */2 */5 * *"
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})

	Context("no match", func() {
		BeforeEach(func() {
			expression = "0 0 1 1 0"
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not match", func() {
			Expect(matches).To(BeFalse())
		})
	})

	Context("invalid number of fields", func() {
		BeforeEach(func() {
			expression = "* * *"
		})

		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("invalid range", func() {
		BeforeEach(func() {
			expression = "60-70 * * * *"
		})

		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("invalid step", func() {
		BeforeEach(func() {
			expression = "*/0 * * * *"
		})

		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("invalid character", func() {
		BeforeEach(func() {
			expression = "a * * * *"
		})

		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("last day of the month", func() {
		BeforeEach(func() {
			expression = "59 23 31 5 3"
			date = time.Date(2023, 5, 31, 23, 59, 0, 0, time.UTC)
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})

	Context("leap year", func() {
		BeforeEach(func() {
			expression = "0 12 29 2 4"
			date = time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("matches", func() {
			Expect(matches).To(BeTrue())
		})
	})
})
