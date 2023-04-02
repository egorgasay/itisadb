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
      } else if (command == "справка") {
        output.innerHTML = "set ключ значение сервер(необязательно) - Устанавливает значение для ключа.<br>сервер > 0 - Сохранить на определенный сервер.<br> сервер = 0 (по умолчанию) - Автоматическое сохранение на менее загруженный сервер. <br>сервер = -1 - Прямое сохранение в базу данных на жестком диске.<br> сервер = -2 - Сохранение во всех экземплярах базы данных.<br>  сервер = -3 - Сохранение во всех экземплярах и базе данных на жестком диске. <br><br>get key ключ сервер(необязательно) - Получает значение из хранилища.<br>сервер > 0 - Поиск на определенном сервере. (скорость: быстрая)<br>сервер = 0 (по умолчанию) - глубокий поиск. (скорость: медленная)<br>сервер = -1 - поиск по базе данных на жестком диске. (скорость: средняя)<br><br>history - История действий пользователя.<br>servers - Список активных серверов со статистикой.";
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