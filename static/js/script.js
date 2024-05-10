window.onload = function() {
    var name = document.getElementById("name");
    name.focus();
    name.addEventListener("keyup", function(e) {
        if (e.keyCode === 13 && name.value != "") {
            let responseText = "";
            var split = this.value.split(" ");
            var command = split[0];
            if (command === "help") {
                if (split.length === 1) {
                    setResponse(name, "Usage: help COMMAND<br>" +
                        "<br>" + "Available commands: <br>" +
                        "<br>Basic commands: <br>" +
                        "- SET - Sets the value to the storage, <br>" +
                        "- GET - Gets the value from the storage, <br>" +
                        "- DEL - Deletes the key-value pair from the storage. <br>" +
                        "<br>Data manipulation: <br>" +
                        "- NEW OBJECT - Creates an object with the specified name, <br>" +
                        "- DELETE OBJECT - Deletes the object with the specified name. <br>" +
                        "<br>Basic object commands: <br>" +
                        "- SETO - Sets the value to the object, <br>" +
                        "- GETO - Gets the value from the object, <br>" +
                        "- DELO - Deletes the object key. <br>" +
                        "<br>Object manipulation: <br>" +
                        "- MARSHAL OBJECT - Displays the object as JSON, <br>" +
                        "- ATTACH - Attaches the src object to the dst. <br>" +
                        "<br>Access control: <br>" +
                        "- NEW USER - Creates a new user, <br>" +
                        "- CHANGE PASSWORD - Changes the password of the user, " +
                        "- CHANGE LEVEL - Changes the level of the user, " +
                        "- DELETE USER - Deletes the user. " +
                        "<br>Advanced: <br>" +
                        "- HISTORY - History of user actions, <br>" +
                        "- SERVERS - List of active balancer with stats. ");
                } else {
                    switch (split[1].toLowerCase()) {
                        case "set":
                            window.location.href = '/doc/base-commands.html#set';
                            break;
                        case "get":
                            window.location.href = '/doc/base-commands.html#get';
                            break;
                        case "del":
                            window.location.href = '/doc/base-commands.html#del';
                            break;
                        case "seto":
                            window.location.href = '/doc/object-commands.html#seto';
                            break;
                        case "geto":
                            window.location.href = '/doc/object-commands.html#geto';
                            break;
                        case "delo":
                            window.location.href = '/doc/object-commands.html#delo';
                            break;
                        case "new":
                            switch (split[2]) {
                                case "object":
                                    window.location.href = '/doc/object-managment.html#new-object';
                                    break;
                            }
                            break;
                        case "delete":
                            switch (split[2]) {
                                case "object":
                                    window.location.href = '/doc/object-managment.html#delete-object';
                                    break;
                            }
                            break;
                        case "marshal":
                            switch (split[2]) {
                                case "object":
                                    window.location.href = '/doc/advanced-object-commands.html#marshal-object';
                                    break;
                            }
                            break;
                        case "attach":
                            window.location.href = '/doc/advanced-object-commands.html#attach';
                            break;
                    }
                }
            } else if (command === "справка") {
                responseText = "set ключ значение сервер(необязательно) - Устанавливает значение для ключа.<br>сервер > 0 - Сохранить на определенный сервер.<br> сервер = 0 (по умолчанию) - Автоматическое сохранение на менее загруженный сервер. <br>сервер = -1 - Прямое сохранение в базу данных на жестком диске.<br> сервер = -2 - Сохранение во всех экземплярах базы данных.<br>  сервер = -3 - Сохранение во всех экземплярах и базе данных на жестком диске. <br><br>get key ключ сервер(необязательно) - Получает значение из хранилища.<br>сервер > 0 - Поиск на определенном сервере. (скорость: быстрая)<br>сервер = 0 (по умолчанию) - глубокий поиск. (скорость: медленная)<br>сервер = -1 - поиск по базе данных на жестком диске. (скорость: средняя)<br><br>" +
                    "new_object name - Создает объект с указанным именем.<br>object name set attr value - Устанавливает значение атрибута object.<br>object name get attr - Получает значение атрибута object.<br>show_object name - Отображает объект в виде карты.<br>attach dst src - Прикрепляет объект src к dst.<br><br>" +
                    "history - История действий пользователя.<br>balancer - Список активных серверов со статистикой.";
            } else if (command === "history") {
                fetch("/history").then(response => response.json()).then(json => setResponse(name, json.text));
            } else if (command === "servers") {
                fetch("/servers").then(response => response.json()).then(json => setResponse(name, json.text));
            } else if (command === "exit") {
                console.log('Redirecting to /exit');
                document.cookie.split(";").forEach(function(c) {
                    document.cookie = c.trim().split("=")[0] + '=;expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/';
                });
                window.location.href = '/exit';
            } else if (command === "doc") {
                console.log('Redirecting to /doc');
                window.location.href = '/doc/';
            } else {
                fetch("/act?action=" + encodeURIComponent(this.value)).then(response => response.json()).then(json => setResponse(name, json.text));
            }
        }
    });
};

function setResponse(command, text) {
    var output = document.querySelector(".message");

    output.innerHTML = "# " + command.value + "<br>" + text + "<br><br>" + output.innerHTML;
}