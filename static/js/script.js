window.onload = function () {
  var name = document.getElementById("name");
  var output = document.querySelector(".message");
  name.focus();
  name.addEventListener("keyup", function (e) {
    if (e.keyCode == 13 && name.value != "") {
      var command = this.value;
      this.value = "";
      if (command == "help") {
        output.innerHTML = "set key value server(optional) - Sets the value to the storage.<br>server > 0 - Save to exact server.<br> server = 0 (default) - Automatic saving to a less loaded server. <br>server = -1 - Direct saving to the database. <br> server = -2 - Saving in all instances.<br> server = -3 - Saving in all instances and DB.  <br><br>get key server(optional) - Gets the value from the storage.<br>server > 0 - Search on a specific server. (speed: fast)<br>server = 0 (default) - Deep search. (speed: slow)<br>server = -1 - DB search. (speed: medium)<br><br>history - History of user actions.<br>servers - List of active servers with stats.";
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