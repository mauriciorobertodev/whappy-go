package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upload Inputs", func() {
	Describe("UpdateUploadMetadata Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateUploadMetadata{
				FileID:   "550e8400-e29b-41d4-a716-446655440000",
				Metadata: file.Metadata{},
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for invalid UUID", func() {
			inp := &input.UpdateUploadMetadata{
				FileID:   "invalid-uuid",
				Metadata: file.Metadata{},
			}
			Expect(inp.Validate()).To(Equal(file.ErrInvalidFileID))
		})
	})

	Describe("ListUploads Input", func() {
		It("should normalize limit to default when not set", func() {
			inp := &input.ListUploads{}
			inp.Normalize()
			Expect(inp.Limit).To(Equal(20))
		})

		It("should normalize limit to max when exceeding max", func() {
			inp := &input.ListUploads{Limit: 150}
			inp.Normalize()
			Expect(inp.Limit).To(Equal(100))
		})

		It("should keep limit as is when within range", func() {
			inp := &input.ListUploads{Limit: 50}
			inp.Normalize()
			Expect(inp.Limit).To(Equal(50))
		})
	})

	Describe("GetUpload Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetUpload{
				FileID: "550e8400-e29b-41d4-a716-446655440000",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for invalid UUID", func() {
			inp := &input.GetUpload{
				FileID: "invalid-uuid",
			}
			Expect(inp.Validate()).To(Equal(file.ErrInvalidFileID))
		})
	})

	Describe("DeleteUpload Input", func() {
		It("should validate successfully", func() {
			inp := &input.DeleteUpload{
				FileID: "550e8400-e29b-41d4-a716-446655440000",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation for invalid UUID", func() {
			inp := &input.DeleteUpload{
				FileID: "invalid-uuid",
			}
			Expect(inp.Validate()).To(Equal(file.ErrInvalidFileID))
		})
	})
})
