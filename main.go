package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	userService := newUserService()

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Server is up and running!"))
		})
		r.Post("/create", userService.CreateUser)
		r.Get("/{id}", userService.GetUser)
		r.Post("/make_friends", userService.MakeFriends)
		r.Delete("/user", userService.DeleteUser)
		r.Get("/friends/{id}", userService.GetFriends)
		r.Put("/{id}", userService.UpdateAge)
	})

	http.ListenAndServe(":3333", r)
}

func (us *userService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	if err := render.Bind(r, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := us.storeUser(user)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, id)
}

func (us *userService) UpdateAge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newAge := &NewAge{}
	if err := render.Bind(r, newAge); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !us.updateAge(id, newAge.NewAge) {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	render.PlainText(w, r, "Возраст пользователя успешно обновлён")
}

func (us *userService) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := us.getUser(id)
	if user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	render.Render(w, r, user)
}

func (us *userService) GetFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	friends, ok := us.getFriends(id)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	render.JSON(w, r, friends)
}

func (us *userService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	target := &TargetID{}
	if err := render.Bind(r, target); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name, ok := us.deleteUser(target.TargetID)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	render.PlainText(w, r, name)
}

func (us *userService) MakeFriends(w http.ResponseWriter, r *http.Request) {
	friends := &Friends{}
	if err := render.Bind(r, friends); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sName, tName, ok := us.makeFriends(friends)
	if !ok {
		http.Error(w, "Друг не найден по id", http.StatusBadRequest)
		return
	}
	render.PlainText(w, r, fmt.Sprintf("%s и %s теперь друзья", sName, tName))
}
