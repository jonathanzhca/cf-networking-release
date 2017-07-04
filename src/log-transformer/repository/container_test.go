package repository_test

import (
	"lib/datastore"
	"lib/fakes"
	"log-transformer/repository"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	var (
		repo      *repository.ContainerRepo
		fakeStore *fakes.Datastore
	)
	BeforeEach(func() {
		fakeStore = &fakes.Datastore{}

		repo = &repository.ContainerRepo{
			Store: fakeStore,
		}

		fakeStore.ReadAllReturns(map[string]datastore.Container{
			"handle-1": {
				Handle: "handle-1",
				IP:     "ip-1",
				Metadata: map[string]interface{}{
					"app_id":   "app-1",
					"space_id": "space-1",
					"org_id":   "org-1",
				},
			},
			"handle-2": {
				Handle:   "handle-2",
				IP:       "ip-2",
				Metadata: map[string]interface{}{},
			},
		}, nil)
	})

	Describe("GetByIP", func() {
		It("looks up the container from the store", func() {
			container, err := repo.GetByIP("ip-1")
			Expect(err).NotTo(HaveOccurred())

			Expect(container).To(Equal(repository.Container{
				Handle:  "handle-1",
				AppID:   "app-1",
				SpaceID: "space-1",
				OrgID:   "org-1",
			}))
		})

		It("looks up the container from the store", func() {
			container, err := repo.GetByIP("ip-1")
			Expect(err).NotTo(HaveOccurred())

			Expect(container).To(Equal(repository.Container{
				Handle:  "handle-1",
				AppID:   "app-1",
				SpaceID: "space-1",
				OrgID:   "org-1",
			}))
		})
	})
})
