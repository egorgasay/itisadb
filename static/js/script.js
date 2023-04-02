window.onload = function () {
  var name = document.getElementById("name");
  var output = document.querySelector(".message");
  name.focus();
  name.addEventListener("keyup", function (e) {
    if (e.keyCode == 13 && name.value != "") {
      var command = this.value;
      this.value = "";
      if (command == "help") {
        output.innerHTML = "set key value - Sets the value to the storage.<br>get key server(optional) - Gets the value from the storage.<br>server > 0 - Search on a specific server. (speed: fast)<br>server = 0 - DB search. (speed: medium)<br>server = -1 (default) - Deep search. (speed: slow)<br><br>history - history of user actions<br>servers - list of active servers with stats";
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