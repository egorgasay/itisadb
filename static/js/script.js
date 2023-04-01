window.onload = function () {
  var name = document.getElementById("name");
  var output = document.querySelector(".message");
  name.focus();
  name.addEventListener("keyup", function (e) {
    if (e.keyCode == 13 && name.value != "") {
      var command = this.value;
      this.value = "";
      if (command == "help") {
        output.innerHTML = "set key value <br>get key<br>history - history of user actions<br>servers - list of active servers with stats";
      } else if (command == "history") {
        fetch("/history").then(response => response.json()).then(json => output.innerHTML = json.text);
      } else if (command == "servers") {
        fetch("/servers").then(response => response.json()).then(json => output.innerHTML = json.text);
      }else {
        fetch("/act?action="+command).then(response => response.json()).then(json => output.innerHTML = json.text);
      }
    }
  });
};