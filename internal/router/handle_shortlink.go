package router

import (
	"GoShort/database"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
)

func getBrowserInfo(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Chrome") {
		return "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		return "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		return "Safari"
	}
	return "Unknown"
}

func getCountryInfo(r *http.Request) string {
	// TODO: integrate with a service later (preferably GeoIP2)
	return "Unknown"
}

func handleShortlink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortlink := vars["shortlink"]
	redirectShortlink, err := database.GetShortlinkByShortURL(shortlink)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	redirectURL := redirectShortlink.LongURL + "?utm_source=goshort&utm_medium=redirect&utm_campaign=" + shortlink

	isHTTP := strings.HasPrefix(redirectShortlink.LongURL, "http://")

	browser := getBrowserInfo(r)
	country := getCountryInfo(r)

	err = database.RecordClickWithDetails(redirectShortlink.ID, browser, country)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	warningTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Redirecting...</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            color: #e0e0e0;
            background-color: #121212;
            margin: 0;
            padding: 0;
            text-align: center;
        }
        h1 {
            color: #f5f5f5;
        }
        a {
            color: #bb86fc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        footer {
            position: fixed;
            bottom: 0;
            width: 100%;
            background-color: #1f1f1f;
            padding: 10px;
            color: #888;
        }
        .warning {
            background-color: #ffeb3b;
            color: #000;
            padding: 10px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            color: #fff;
            background-color: #bb86fc;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #9b72ea;
        }
    </style>
    <script>
        function redirect() {
            const redirectUrl = "{{.RedirectUrl}}";
            window.location.href = redirectUrl;
        }
    </script>
</head>
<body>
    <h1>Redirecting...</h1>
    {{if .IsHTTP}}
        <div class="warning">
            <p><strong>Warning:</strong> You are being redirected to a non-secure HTTP site. Proceed with caution.</p>
			<p>Do NOT enter any personal information on that website. Such as passwords or bank details!</p>
            <p>If you wish to continue, please <a href="#" class="button" onclick="redirect()">Click Here</a>.</p>
        </div>
    {{else}}
        <script>
            setTimeout(redirect, 3000);
        </script>
    {{end}}
    <p>If you are not redirected automatically, <a href="{{.RedirectUrl}}">click here</a>.</p>
    <footer>
        <p>Powered by <a href="https://github.com/tkbstudios/goshort-server" target="_blank">GoShort! by TKB Studios</a> | Not liable for hosted links.</p>
    </footer>
</body>
</html>
`

	t, err := template.New("redirect").Parse(warningTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		RedirectUrl string
		IsHTTP      bool
	}{
		RedirectUrl: redirectURL,
		IsHTTP:      isHTTP,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
