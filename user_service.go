package main

import (
	"strconv"
)

type userService struct {
	usersMap  map[string]*User
	idCounter int
}

func newUserService() *userService {
	return &userService{usersMap: make(map[string]*User)}
}

func (us *userService) storeUser(user *User) string {
	us.idCounter++
	user.ID = strconv.Itoa(us.idCounter)
	us.usersMap[user.ID] = user
	return user.ID
}

func (us *userService) getUser(id string) *User {
	user, ok := us.usersMap[id]
	if ok {
		return user
	}
	return nil
}

func (us *userService) deleteUser(targetID string) (string, bool) {
	userToDelete, ok := us.usersMap[targetID]
	if ok {
		for _, friendID := range userToDelete.Friends {
			us.deleteFromFriends(friendID, userToDelete.ID)
		}

		delete(us.usersMap, targetID)
		return userToDelete.Name, true
	}
	return "", false
}

func (us *userService) deleteFromFriends(srcUserID string, toDeleteID string) {
	u := us.getUser(srcUserID)
	deleteIdx := -1
	for i, userID := range u.Friends {
		if userID == toDeleteID {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		u.Friends = append(u.Friends[:deleteIdx], u.Friends[deleteIdx+1:]...)
	}
}

func (us *userService) makeFriends(friends *Friends) (string, string, bool) {
	srcUser := us.getUser(friends.SourceID)
	targetUser := us.getUser(friends.TargetID)
	if srcUser == nil || targetUser == nil {
		return "", "", false
	}

	for _, friendID := range srcUser.Friends {
		if friendID == targetUser.ID {
			return srcUser.Name, targetUser.Name, true
		}
	}

	srcUser.Friends = append(srcUser.Friends, friends.TargetID)
	targetUser.Friends = append(targetUser.Friends, friends.SourceID)
	return srcUser.Name, targetUser.Name, true
}

func (us *userService) getFriends(userID string) ([]*User, bool) {
	user := us.getUser(userID)
	if user == nil {
		return nil, false
	}
	var friends []*User
	for _, friendID := range user.Friends {
		friends = append(friends, us.getUser(friendID))
	}
	return friends, true
}

func (us *userService) updateAge(userID string, age string) bool {
	user := us.getUser(userID)
	if user == nil {
		return false
	}
	user.Age = age
	return true
}
