package sqlite

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (c *Config) SetContainerOwner(username, name, containerid string) error {
	var user User
	c.DB.Where("name = ?", username).First(&user)
	if user.ID == 0 {
		return nil // FIXME
	}

	co := ContainerOwner{
		ContainerID: containerid,
		User:        user,
	}

	if len(name) > 1 { // container name must be at least 2 characters
		co.Name = name

	}
	// FIXME Check for error
	c.DB.Model(&ContainerOwner{}).Create(&co)
	return nil
}

func (c *Config) RemoveContainerOwner(username, name, containerid string) error {
	var co ContainerOwner
	var u User

	c.DB.Where("name = ?", username).First(&u)
	if u.ID == 0 {
		return nil
	}

	if len(name) > 1 {
		result := c.DB.Debug().Where("name = ? AND user_id = ?", name, u.ID).Delete(&co)
		log.Debug(result.RowsAffected)
	}

	result := c.DB.Debug().Where("container_id = ? AND user_id = ?", containerid, u.ID).Delete(&co)
	log.Debug(result.RowsAffected)

	prefix := fmt.Sprintf("%s%%", containerid)
	result = c.DB.Debug().Where("container_id LIKE ? AND user_id = ?", prefix, u.ID).Delete(&co)
	log.Debug(result.RowsAffected)

	return nil
}

func (c *Config) IsContainerOwner(username, containerid string) bool {
	var co ContainerOwner
	var u User
	var cnt int

	c.DB.Where("name = ?", username).First(&u)
	if u.ID == 0 {
		return false
	}

	c.DB.Model(&co).Where("container_id = ? AND user_id = ?", containerid, u.ID).Count(&cnt)
	if cnt == 1 {
		return true
	}

	c.DB.Model(&co).Where("name = ? AND user_id = ?", containerid, u.ID).Count(&cnt)
	if cnt == 1 {
		return true
	}

	prefix := fmt.Sprintf("%s%%", containerid)
	prfm := false
	var cop []ContainerOwner
	c.DB.Where("container_id LIKE ?", prefix).Find(&cop)

	for _, p := range cop {
		if p.UserID != u.ID {
			return false
		}
		prfm = true
	}

	return prfm
}
