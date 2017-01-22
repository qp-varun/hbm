package sqlite

import (
	"github.com/kassisol/hbm/storage/driver"
)

func (c *Config) AddResource(name, rtype, value, options string) {
	c.DB.Create(&Resource{
		Name:   name,
		Type:   rtype,
		Value:  value,
		Option: options,
	})
}

func (c *Config) RemoveResource(name string) {
	c.DB.Where("name = ?", name).Delete(Resource{})
}

func (c *Config) ListResources() []driver.ResourceResult {
	var resources []Resource
	var result []driver.ResourceResult

	c.DB.Find(&resources)

	for _, res := range resources {
		result = append(result, driver.ResourceResult{Name: res.Name, Type: res.Type, Value: res.Value, Option: res.Option})
	}

	return result
}

func (c *Config) FindResource(name string) bool {
	var count int64

	c.DB.Model(&Resource{}).Where("name = ?", name).Count(&count)

	if count == 1 {
		return true
	}

	return false
}

func (c *Config) CountResource(rtype string) int {
	var count int64

	if rtype == "all" {
		c.DB.Model(&Resource{}).Count(&count)
	} else {
		c.DB.Model(&Resource{}).Where("type = ?", rtype).Count(&count)
	}

	return int(count)
}

func (c *Config) AddResourceToCollection(collection, resource string) {
	col := Collection{}
	res := Resource{}

	c.DB.Where("name = ?", resource).Find(&res)
	c.DB.Where("name = ?", collection).Find(&col)

	c.DB.Model(&col).Association("Resources").Append(&res)
}

func (c *Config) RemoveResourceFromCollection(collection, resource string) {
	col := Collection{}
	res := Resource{}

	c.DB.Where("name = ?", resource).Find(&res)
	c.DB.Where("name = ?", collection).Find(&col)

	c.DB.Model(&col).Association("Resources").Delete(&res)
}