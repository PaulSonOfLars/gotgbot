package main

import (
	"fmt"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func index(webappURL string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Telegram webapp</title>
	<script src="https://telegram.org/js/telegram-web-app.js"></script>
</head>
<body>
<div id="main">
	Hey there :)
	This is a VERY basic example of a webapp.

	All it does is read the info from telegram, and validate it.

	</br>
	<p id="name"></p>
	</br>
	<p id="id"></p>
	</br>
	<p id="valid">unchecked</p>
</div>

<script>
Telegram.WebApp.ready()

document.getElementById("name").innerHTML = "your name is: " + Telegram.WebApp.initDataUnsafe.user.first_name
document.getElementById("id").innerHTML = "your id is: " + Telegram.WebApp.initDataUnsafe.user.id

fetch("%s/validate?"+Telegram.WebApp.initData).then(function(response) {
	return response.text();
}).then(function(text) {
  document.getElementById("valid").innerHTML = "result: " + text;
}).catch(function() {
  console.log("Booo");
});
</script>
</body>
</html>
`, webappURL)))
	}
}

func validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if ext.ValidateWebAppQuery(request.URL.Query(), token) {
			writer.Write([]byte("validation success; user is authenticated."))
		} else {
			writer.Write([]byte("validation failed; data cannot be trusted."))
		}
	}
}
