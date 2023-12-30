window.onload = function () {
  var name = document.getElementById("name");
  var output = document.querySelector(".message");
  name.focus();
  name.addEventListener("keyup", function (e) {
    if (e.keyCode == 13 && name.value != "") {
      var split = this.value.split(" ");
      var command = split[0];
      if (command === "help") {
        if (split.length === 1) {
          output.innerHTML = "Usage: help COMMAND<br>" +
              "<br>" + "Available commands: <br>" +
              "<br>Basic commands: <br>" +
              "- SET - Sets the value to the storage, <br>" +
              "- GET - Gets the value from the storage, <br>" +
              "- DEL - Deletes the key-value pair from the storage, <br>" +
              "<br>Data manipulation: <br>" +
              "- NEW OBJECT - Creates an object with the specified name, <br>" +
              "- DELETE OBJECT - Deletes the object with the specified name, <br>" +
              "<br>Basic object commands: <br>" +
              "- SETO - Sets the value to the object, <br>" +
              "- GETO - Gets the value from the object, <br>" +
              "- DELO - Deletes the object key, <br>" +
              "<br>Object manipulation: <br>" +
              "- MARSHAL OBJECT - Displays the object as JSON, <br>" +
              "- ATTACH - Attaches the src object to the dst, <br>" +
              "<br>Advanced: <br>" +
              "- HISTORY - History of user actions, <br>" +
              "- SERVERS - List of active servers with stats. ";
        } else {
          switch (split[1].toLowerCase()) {
            case "set":
              output.innerHTML = "SET key \"value\" [ MODE - NX | RO | XX ] [ LEVEL - R | S ] [ SERVER - [0-9]+ ]  <br>" +
                  "<br>" +
                  "MODE - Defines the mode of the operation. <br>" +
                  "- `NX` - If the key already exists, it won't be overwritten. <br>" +
                  "- `RO` - If the key already exists, an error will be returned.<br>" +
                  "- `XX` - If the key doesn't exist, it won't be created.<br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission. <br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation<br>" +
                  "By default - NO encryption, NO ACL validation<br>" +
                  "<br>" +
                  "SERVER - Defines server number to use. <br>" +
                  "- Automaticly saving to a less loaded server by default.";
              break;
            case "get":
              output.innerHTML = "GET key [ FROM - [0-9]+ ] [ LEVEL - D | R | L ]  <br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation<br>" +
                  "<br>" +
                  "FROM - Defines server number to use.<br>" +
                  "- `> 0` - Search on a specific server (speed: fast).  <br>" +
                  "- `= 0` (default) - Deep search (speed: slow). ";
              break;
            case "del":
              output.innerHTML = "DEL key [ FROM - [0-9]+ ] [ LEVEL - D | R | L ]  <br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation<br>" +
                  "<br>" +
                  "FROM - Defines server number to use.<br>" +
                  "- `> 0` - Search on a specific server (speed: fast).  <br>" +
                  "- `= 0` (default) - Deep search (speed: slow). ";
              break;
            case "seto":
              output.innerHTML = "SETO name key \"value\" [ LEVEL - D | R | L ]<br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation";
              break;
            case "geto":
              output.innerHTML = "GETO name key [ LEVEL - D | R | L ]<br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation";
              break;
            case "delo":
              output.innerHTML = "DELO name key [ LEVEL - D | R | L ]<br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation";
              break;
            case "new":
              switch (split[2]) {
                case "object":
                  output.innerHTML = "NEW OBJECT name [ ON - [0-9]+ ] [ LEVEL - D | R | L ]<br>" +
                      "<br>" +
                      "LEVEL - Defines the level of permission.<br>" +
                      "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                      "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                      "- `S` (Secret) - encryption, ACL validation<br>" +
                      "<br>" +
                      "ON - Defines server number to use. <br>" +
                      "- Automaticly saving to a less loaded server by default.";
                  break;
              }
              break;
            case "delete":
              switch (split[2]) {
                case "object":
                  output.innerHTML = "DELETE OBJECT name [ ON - [0-9]+ ] [ LEVEL - D | R | L ]<br>" +
                      "<br>" +
                      "LEVEL - Defines the level of permission.<br>" +
                      "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                      "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                      "- `S` (Secret) - encryption, ACL validation<br>" +
                      "<br>" +
                      "ON - Defines server number to use. <br>" +
                      "- Automaticly saving to a less loaded server by default.";
                  break;
              }
              break;
            case "marshal":
              switch (split[2]) {
                case "object":
                  output.innerHTML = "MARSHAL OBJECT name [ LEVEL - D | R | L ]<br>" +
                      "<br>" +
                      "LEVEL - Defines the level of permission.<br>" +
                      "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                      "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                      "- `S` (Secret) - encryption, ACL validation";
                  break;
              }
              break;
            case "attach":
              output.innerHTML = "ATTACH dst [ LEVEL - D | R | L ] src [ LEVEL - D | R | L ]<br>" +
                  "<br>" +
                  "LEVEL - Defines the level of permission.<br>" +
                  "- `D` (Default) - NO encryption, NO ACL validation<br>" +
                  "- `R` (Restricted) - NO encryption, ACL validation<br>" +
                  "- `S` (Secret) - encryption, ACL validation";
              break;
          }
        }
      } else if (command == "справка") {
        output.innerHTML = "set ключ значение сервер(необязательно) - Устанавливает значение для ключа.<br>сервер > 0 - Сохранить на определенный сервер.<br> сервер = 0 (по умолчанию) - Автоматическое сохранение на менее загруженный сервер. <br>сервер = -1 - Прямое сохранение в базу данных на жестком диске.<br> сервер = -2 - Сохранение во всех экземплярах базы данных.<br>  сервер = -3 - Сохранение во всех экземплярах и базе данных на жестком диске. <br><br>get key ключ сервер(необязательно) - Получает значение из хранилища.<br>сервер > 0 - Поиск на определенном сервере. (скорость: быстрая)<br>сервер = 0 (по умолчанию) - глубокий поиск. (скорость: медленная)<br>сервер = -1 - поиск по базе данных на жестком диске. (скорость: средняя)<br><br>"
        + "new_object name - Создает объект с указанным именем.<br>object name set attr value - Устанавливает значение атрибута object.<br>object name get attr - Получает значение атрибута object.<br>show_object name - Отображает объект в виде карты.<br>attach dst src - Прикрепляет объект src к dst.<br><br>"
         + "history - История действий пользователя.<br>servers - Список активных серверов со статистикой.";
         } else if (command == "history") {
        fetch("/history").then(response => response.json()).then(json => output.innerHTML = json.text);
      } else if (command == "servers") {
        fetch("/servers").then(response => response.json()).then(json => output.innerHTML = json.text);
      }else {
        fetch("/act?action="+encodeURIComponent(this.value)).then(response => response.json()).then(json => output.innerHTML = json.text);
      }
    }
  });
};