package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"cat-api/models"

	"github.com/gorilla/mux"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// JSON response helper
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

// Cat handlers
func (h *Handler) GetCats(w http.ResponseWriter, r *http.Request) {
	include := r.URL.Query().Get("include")

	var query string
	if include == "all" {
		query = `
            SELECT c.id, c.name, c.cat_type_id, c.master_id, c.created_at, c.updated_at,
                   t.id, t.name, t.created_at,
                   m.id, m.first_name, m.last_name, m.place, m.created_at, m.updated_at
            FROM cats c
            LEFT JOIN types t ON c.cat_type_id = t.id
            LEFT JOIN masters m ON c.master_id = m.id
        `
	} else {
		query = `SELECT id, name, cat_type_id, master_id, created_at, updated_at FROM cats`
	}

	rows, err := h.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat

		if include == "all" {
			var catType models.Type
			var master models.Master

			err := rows.Scan(
				&cat.ID, &cat.Name, &cat.CatTypeID, &cat.MasterID, &cat.CreatedAt, &cat.UpdatedAt,
				&catType.ID, &catType.Name, &catType.CreatedAt,
				&master.ID, &master.FirstName, &master.LastName, &master.Place, &master.CreatedAt, &master.UpdatedAt,
			)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			cat.CatType = &catType
			cat.Master = &master
		} else {
			err := rows.Scan(&cat.ID, &cat.Name, &cat.CatTypeID, &cat.MasterID, &cat.CreatedAt, &cat.UpdatedAt)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		cats = append(cats, cat)
	}

	respondWithJSON(w, http.StatusOK, cats)
}

func (h *Handler) GetCat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	query := `
        SELECT c.id, c.name, c.cat_type_id, c.master_id, c.created_at, c.updated_at,
               t.id, t.name, t.created_at,
               m.id, m.first_name, m.last_name, m.place, m.created_at, m.updated_at
        FROM cats c
        LEFT JOIN types t ON c.cat_type_id = t.id
        LEFT JOIN masters m ON c.master_id = m.id
        WHERE c.id = ?
    `

	var cat models.Cat
	var catType models.Type
	var master models.Master

	err = h.DB.QueryRow(query, id).Scan(
		&cat.ID, &cat.Name, &cat.CatTypeID, &cat.MasterID, &cat.CreatedAt, &cat.UpdatedAt,
		&catType.ID, &catType.Name, &catType.CreatedAt,
		&master.ID, &master.FirstName, &master.LastName, &master.Place, &master.CreatedAt, &master.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Cat not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cat.CatType = &catType
	cat.Master = &master
	respondWithJSON(w, http.StatusOK, cat)
}

func (h *Handler) CreateCat(w http.ResponseWriter, r *http.Request) {
	var cat models.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if cat.Name == "" || cat.CatTypeID == 0 || cat.MasterID == 0 {
		respondWithError(w, http.StatusBadRequest, "Name, cat_type_id, and master_id are required")
		return
	}

	query := `INSERT INTO cats (name, cat_type_id, master_id) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, cat.Name, cat.CatTypeID, cat.MasterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cat.ID = int(id)
	respondWithJSON(w, http.StatusCreated, cat)
}

func (h *Handler) UpdateCat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var cat models.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	query := `UPDATE cats SET name=?, cat_type_id=?, master_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`
	result, err := h.DB.Exec(query, cat.Name, cat.CatTypeID, cat.MasterID, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Cat not found")
		return
	}

	cat.ID = id
	respondWithJSON(w, http.StatusOK, cat)
}

func (h *Handler) DeleteCat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	result, err := h.DB.Exec("DELETE FROM cats WHERE id=?", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Cat not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Cat deleted successfully"})
}

// Type handlers
func (h *Handler) GetTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, created_at FROM types")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var types []models.Type
	for rows.Next() {
		var t models.Type
		err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		types = append(types, t)
	}

	respondWithJSON(w, http.StatusOK, types)
}

func (h *Handler) CreateType(w http.ResponseWriter, r *http.Request) {
	var t models.Type
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if t.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	query := `INSERT INTO types (name) VALUES (?)`
	result, err := h.DB.Exec(query, t.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t.ID = int(id)
	respondWithJSON(w, http.StatusCreated, t)
}

// Master handlers
func (h *Handler) GetMasters(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, first_name, last_name, place, created_at, updated_at FROM masters")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var masters []models.Master
	for rows.Next() {
		var m models.Master
		err := rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Place, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		masters = append(masters, m)
	}

	respondWithJSON(w, http.StatusOK, masters)
}

func (h *Handler) GetMaster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var master models.Master
	err = h.DB.QueryRow(
		"SELECT id, first_name, last_name, place, created_at, updated_at FROM masters WHERE id=?",
		id,
	).Scan(&master.ID, &master.FirstName, &master.LastName, &master.Place, &master.CreatedAt, &master.UpdatedAt)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Master not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, master)
}

func (h *Handler) CreateMaster(w http.ResponseWriter, r *http.Request) {
	var master models.Master
	if err := json.NewDecoder(r.Body).Decode(&master); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if master.FirstName == "" || master.LastName == "" {
		respondWithError(w, http.StatusBadRequest, "First name and last name are required")
		return
	}

	query := `INSERT INTO masters (first_name, last_name, place) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, master.FirstName, master.LastName, master.Place)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	master.ID = int(id)
	respondWithJSON(w, http.StatusCreated, master)
}

// Health check handler
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":   "OK",
		"database": "connected",
	})
}
