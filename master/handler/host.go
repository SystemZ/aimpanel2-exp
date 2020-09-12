package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/metric"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/response"
	"gitlab.com/systemz/aimpanel2/master/service/host"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

// @Router /host [get]
// @Summary List
// @Tags Host
// @Description List Hosts linked to the current signed-in account
// @Accept json
// @Produce json
// @Success 200 {object} response.HostList
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostList(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)
	hosts, err := model.GetHostsByUser(user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.HostList{Hosts: hosts})
}

// @Router /host/{id} [get]
// @Summary Details
// @Tags Host
// @Description Get details about Host with selected ID linked to the current signed-in account
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.Host
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	h, err := model.GetHostById(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if h == nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.HostNotFound, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Host{Host: *h})
}

// @Router /host [post]
// @Summary Create
// @Tags Host
// @Description Create new Host linked to the current signed-in account
// @Accept json
// @Produce json
// @Param host body request.HostCreate true " "
// @Success 200 {object} response.Token
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostCreate(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	data := &request.HostCreate{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	h, errCode := host.Create(data, user.ID)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusInternalServerError, errCode, nil)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.Token{Token: h.Token})
}

// @Router /host/{hostId} [put]
// @Summary Edit
// @Tags Host
// @Description Edit host by selected id
// @Accept json
// @Produce json
// @Param hostId path string true "Host ID"
// @Param host body request.HostCreate true " "
// @Success 200 {object} response.Host
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	h, err := model.GetHostById(hostId)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	data := &request.HostCreate{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	if h.Name != data.Name {
		user := context.Get(r, "user").(model.User)
		err = model.SaveAction(
			task.Message{
				TaskId: task.HOST_NAME_CHANGE,
				HostID: h.ID.Hex(),
			},
			user,
			hostId,
			data.Name,
			h.Name,
		)
		if err != nil {
			lib.ReturnError(w, http.StatusInternalServerError, ecode.DbSave, err)
			return
		}

		h.Name = data.Name
		model.Update(h)
	}

	lib.MustEncode(json.NewEncoder(w), response.Host{Host: *h})
}

// FIXME add URL params
// @Router /host/{id}/metric [get]
// @Summary Metric
// @Tags Host
// @Description Get latest metrics for host with ID, linked to the current signed-in account
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.HostMetrics
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostMetric(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, nil)
		return
	}

	// FIXME more validation for metric params
	// get time between data points
	query := r.URL.Query()
	intervalSStr := query.Get("interval")
	intervalSInt, err := strconv.Atoi(intervalSStr)
	if err != nil {
		// TODO make separate ecode
		lib.ReturnError(w, http.StatusBadRequest, ecode.Unknown, nil)
		return
	}

	// get metric name
	metricName := query.Get("name")
	// FIXME move this to service to make http fat free
	// FIXME use map to make it shorter
	metricId := 0
	switch metricName {
	case "cpu_usage":
		metricId = int(metric.CpuUsage)
	case "cpu_user":
		metricId = int(metric.User)
	case "cpu_system":
		metricId = int(metric.System)
	case "cpu_idle":
		metricId = int(metric.Idle)
	case "cpu_nice":
		metricId = int(metric.Nice)
	case "cpu_guest":
		metricId = int(metric.Guest)
	case "cpu_guest_nice":
		metricId = int(metric.GuestNice)
	case "cpu_steal":
		metricId = int(metric.Steal)
	case "cpu_iowait":
		metricId = int(metric.Iowait)
	case "cpu_irq":
		metricId = int(metric.Irq)
	case "cpu_irq_soft":
		metricId = int(metric.Softirq)
	case "ram_usage":
		metricId = int(metric.RamUsage)
	case "ram_free":
		metricId = int(metric.RamFree)
	case "ram_total":
		metricId = int(metric.RamTotal)
	case "ram_available":
		metricId = int(metric.RamAvailable)
	case "ram_buffers":
		metricId = int(metric.RamBuffers)
	case "ram_cache":
		metricId = int(metric.RamCache)
	case "disk_free":
		metricId = int(metric.DiskFree)
	case "disk_used":
		metricId = int(metric.DiskUsed)
	case "disk_total":
		metricId = int(metric.DiskTotal)
	}
	if metricId == 0 {
		// TODO make separate ecode for incorrect metric name
		lib.ReturnError(w, http.StatusBadRequest, ecode.Unknown, nil)
		return
	}

	lastStr := query.Get("last")
	lastInt, err := strconv.Atoi(lastStr)
	if err != nil {
		// TODO make separate ecode
		lib.ReturnError(w, http.StatusBadRequest, ecode.Unknown, nil)
		return
	}

	now := time.Now()
	from := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year()+1, now.Month(), now.Day(), 23, 59, 0, 0, now.Location())
	if lastInt > 0 {
		nowProcessed := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
		from = nowProcessed.Add(time.Duration(-lastInt) * time.Second)
		to = nowProcessed
		logrus.Debugf("from %v to %v", from, to)
	}

	metrics, err := model.GetTimeSeries(oid, intervalSInt, from, to, metric.Id(metricId))
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	if metrics == nil {
		metrics = make([]model.TimeseriesOutput, 0)
	}

	lib.MustEncode(json.NewEncoder(w), response.HostMetrics{Metrics: metrics})
}

// @Router /host/{id} [delete]
// @Summary Remove
// @Tags Host
// @Description Removes host with all linked game servers
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} response.JsonSuccess
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	oid, _ := primitive.ObjectIDFromHex(params["hostId"])
	user := context.Get(r, "user").(model.User)
	errCode := host.Remove(oid, user)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//TODO: Available for users?
func HostUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, _ := primitive.ObjectIDFromHex(params["hostId"])
	user := context.Get(r, "user").(model.User)
	err := host.Update(oid, user)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.GsUpdate, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{id}/job [post]
// @Summary Create Job
// @Tags Host
// @Description Create new Job linked to the given host
// @Accept json
// @Produce json
// @Param host body request.HostCreateJob true " "
// @Success 200 {object} response.Token
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostCreateJob(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(model.User)

	params := mux.Vars(r)

	data := &request.HostCreateJob{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.JsonDecode, err)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(params["hostId"])
	_, errCode := host.CreateJob(data, user.ID, oid)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusInternalServerError, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Router /host/{id}/job/{job_id} [delete]
// @Summary Remove Job
// @Tags Host
// @Description Removes host job
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Param job_id path string true "Job ID"
// @Success 200 {object} response.JsonSuccess
// @Failure 400 {object} response.JsonError
// @Security ApiKey
func HostJobRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hostId, _ := primitive.ObjectIDFromHex(params["hostId"])
	jobId, _ := primitive.ObjectIDFromHex(params["jobId"])
	errCode := host.RemoveJob(hostId, jobId)
	if errCode != ecode.NoError {
		lib.ReturnError(w, http.StatusBadRequest, errCode, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HostJobList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oid, err := primitive.ObjectIDFromHex(params["hostId"])
	if err != nil {
		lib.ReturnError(w, http.StatusBadRequest, ecode.OidError, err)
		return
	}

	jobs, err := model.GetHostJobsByHostId(oid)
	if err != nil {
		lib.ReturnError(w, http.StatusInternalServerError, ecode.DbError, err)
		return
	}

	lib.MustEncode(json.NewEncoder(w), response.HostJobList{Jobs: jobs})
}
