package server

import (
	"errors"
	"log"
	"net/http"
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"
	"pratbacknd/internal/storage"
	"pratbacknd/internal/utils"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Mux            *chi.Mux
	allowedOrigins string
	storage        storage.Storage
	uuidGen        utils.UUIDGenerator
}

type Config struct {
	AllowedOrigins string
	Storage        storage.Storage
	UUIDGen        utils.UUIDGenerator
}

func New(config Config) (*Server, error) {
	m := chi.NewRouter()
	s := &Server{Mux: m, storage: config.Storage, allowedOrigins: config.AllowedOrigins, uuidGen: config.UUIDGen}

	m.Use(s.enableCORS)

	m.Get("/products", s.Products)
	m.Post("/admin/products", s.CreateProduct)

	m.Get("/categories", s.Categories)
	m.Post("/admin/categories", s.CreateCategory)

	return s, nil
}

func (s *Server) Products(w http.ResponseWriter, r *http.Request) {
	products, err := s.storage.Products()
	if err != nil {
		log.Printf("error - fetching products: %s \n", err)
		s.errorJSON(w, errors.New("error fetching products"), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, products)
}

func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p product.Product
	err := s.readJSON(w, r, &p)

	if err != nil {
		log.Printf("error - building json: %s \n", err)
		s.errorJSON(w, errors.New("error reading product"), http.StatusBadRequest)
		return
	}

	p.ID = s.uuidGen.Generate()

	err = s.storage.CreateProduct(p)
	if err != nil {
		log.Printf("error - storing product: %s \n", err)
		s.errorJSON(w, errors.New("error persisting product"), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, p)
}

func (s *Server) Categories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.storage.Categories()
	if err != nil {
		log.Printf("error - fetching categories: %s \n", err)
		s.errorJSON(w, errors.New("error fetching categories"), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, categories)
}

func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c category.Category
	err := s.readJSON(w, r, &c)

	if err != nil {
		log.Printf("error - building json: %s \n", err)
		s.errorJSON(w, errors.New("error reading category"), http.StatusBadRequest)
		return
	}

	c.ID = s.uuidGen.Generate()

	err = s.storage.CreateCategory(c)
	if err != nil {
		log.Printf("error - storing category: %s \n", err)
		s.errorJSON(w, errors.New("error persisting product"), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, c)
}
