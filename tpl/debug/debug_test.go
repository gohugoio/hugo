package debug

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

type User struct {
	Name  string
	Phone string
	city  string
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetPhone() string {
	return u.Phone
}

func (u *User) getCity() string {
	return u.city
}

func (u *User) GetPhoneAndCity() string {
	return u.Phone + u.city
}

func TestList(t *testing.T) {
	user := User{"a name", "9876543210", "SF"}
	ns := Namespace{}

	t.Run("struct", func(t *testing.T) {
		c := qt.New(t)
		result := ns.List(user)
		c.Assert(len(result), qt.Equals, 4)
		c.Assert(result[0], qt.Equals, "GetName")
		c.Assert(result[1], qt.Equals, "GetPhone")
		c.Assert(result[2], qt.Equals, "Name")
		c.Assert(result[3], qt.Equals, "Phone")
	})

	t.Run("pointer", func(t *testing.T) {
		c := qt.New(t)
		result := ns.List(&user)
		c.Assert(len(result), qt.Equals, 5)
		c.Assert(result[0], qt.Equals, "GetName")
		c.Assert(result[1], qt.Equals, "GetPhone")
		c.Assert(result[2], qt.Equals, "GetPhoneAndCity")
		c.Assert(result[3], qt.Equals, "Name")
		c.Assert(result[4], qt.Equals, "Phone")
	})

	t.Run("map", func(t *testing.T) {
		c := qt.New(t)
		mapTestCase := map[string]string{
			"name":  "a name",
			"phone": "a phone",
		}
		result := ns.List(mapTestCase)
		c.Assert(result[0], qt.Equals, "name")
		c.Assert(result[1], qt.Equals, "phone")
	})

}
