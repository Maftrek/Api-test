package application

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"

	"template/models"
)

// хелф хендлер для кубера
func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	sum := models.GetHashSum()
	if !bytes.Equal(sum, app.hashSum) {
		app.logger.Log("msg", "New Configuration")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *Application) createNews(w http.ResponseWriter, r *http.Request) {
	newsInfo, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.logger.Log("err", err, "body", r.Body)
		w.WriteHeader(http.StatusInternalServerError)
	}

	defer r.Body.Close()
	// инициализируем переменную для получения информации по поставщику
	var news models.Data
	// анмаршалим в структуру
	err = json.Unmarshal(newsInfo, &news)
	if err != nil {
		app.logger.Log("err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp, err := app.svc.CreateNews(news)
	if err != nil {
		app.logger.Log("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		app.logger.Log("err", err)
	}
}

func (app *Application) getNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsID, err := strconv.ParseInt(vars["news_id"], 10, 32)
	if err != nil {
		app.logger.Log("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	news, err := app.svc.GetNews(int32(newsID))
	if err != nil {
		app.logger.Log("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(news)
	if err != nil {
		app.logger.Log("err", err)
	}
}
