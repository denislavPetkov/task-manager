package nlp

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nlp", func() {

	Describe("GetCommand", func() {
		Context("Incorrect input", func() {
			It("should return an empty command when no input", func() {
				command := GetCommand("")
				Expect(command).To(Equal(Command{}))
			})
		})

		Context("Correct input", func() {
			It("should return a correct command when input is correct", func() {
				command := GetCommand("please create new task")
				Expect(command.Command).To(Equal("create"))
			})

			It("should return a correct command when input is correct", func() {
				command := GetCommand("please update the task buy groceries")
				Expect(command).To(Equal(Command{
					Command:   "update the",
					TaskTitle: "buy groceries",
				}))
			})

			It("should return a correct command when input is correct", func() {
				command := GetCommand("please cancel the task repair my computer")
				Expect(command).To(Equal(Command{
					Command:   "cancel the",
					TaskTitle: "repair my computer",
				}))
			})

			It("should return a correct command when input is correct", func() {
				command := GetCommand("please remove the do homework")
				Expect(command).To(Equal(Command{
					Command:   "remove",
					TaskTitle: "do homework",
				}))
			})
		})

	})
})
