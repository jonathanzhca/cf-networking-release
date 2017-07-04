package repository

import "lib/datastore"

type Container struct {
	Handle  string
	AppID   string
	SpaceID string
	OrgID   string
}

type ContainerRepo struct {
	Store datastore.Datastore
}

func (c *ContainerRepo) GetByIP(ip string) (Container, error) {
	containers, err := c.Store.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if container.IP == ip {
			appID, ok := container.Metadata["app_id"].(string)
			if !ok {
				panic("foo")
			}
			spaceID, ok := container.Metadata["space_id"].(string)
			if !ok {
				panic("bar")
			}
			orgID, ok := container.Metadata["org_id"].(string)
			if !ok {
				panic("baz")
			}
			return Container{
				Handle:  container.Handle,
				AppID:   appID,
				SpaceID: spaceID,
				OrgID:   orgID,
			}, nil
		}
	}

	return Container{}, nil
}
