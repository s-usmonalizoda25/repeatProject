package models

import "errors"

type User struct{
	ID int
	Name string
	Age int
}

func (u *User)Validate()error{
	if u.ID<=0{
		return errors.New("invalid user id")
	}
	if u.Age<=0{
		return errors.New("invalid user age")
	}
	if u.Name<=""{
		return errors.New("invalid user name")
	}
	return nil
}

