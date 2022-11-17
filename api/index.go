package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/timotm/tuntihinta-tallentaja/pkg/fetcher"
	"github.com/timotm/tuntihinta-tallentaja/pkg/glue"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token := r.Header.Get("Authorization")
		split := strings.Split(token, " ")
		if len(split) != 2 ||
			split[0] != "Bearer" ||
			split[1] != strings.TrimSpace(os.Getenv("TH_REQUEST_TOKEN")) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		t := time.Now().AddDate(0, 0, 1)
		startTime := fetcher.EntsoeTime(t)

		t = time.Now().AddDate(0, 0, 2)
		endTime := fetcher.EntsoeTime(t)

		fmt.Fprintf(w, "Fetching data for %+v to %+v\n", startTime.String(), endTime.String())
		files := glue.FetchAndUpload(startTime,
			endTime,
			os.Getenv("TH_SECURITY_TOKEN"),
			os.Getenv("TH_AWS_REGION"),
			os.Getenv("TH_AWS_BUCKET_NAME"),
			os.Getenv("TH_AWS_ACCESS_KEY_ID"),
			os.Getenv("TH_AWS_SECRET_ACCESS_KEY"))

		fmt.Fprintf(w, "Uploaded files: %s\n", strings.Join(files, ", "))

	} else {
		http.Error(w, "I don't even", http.StatusMethodNotAllowed)
	}
}
