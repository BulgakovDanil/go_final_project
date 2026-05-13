package api

import (
	"net/http"
	"time"

	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = db.DeletTask(task.ID)
		if err != nil {
			writeJsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeEmptyJson(w, http.StatusOK)
		return
	}

	now := time.Now()

	date, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateDate(task, date)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeEmptyJson(w, http.StatusOK)
}
