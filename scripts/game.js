
window.onload = function () {
    var conn;
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function logMessage(msg) {
        var item = document.createElement("div");
        item.innerHTML = msg;
        appendLog(item);
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        // hide submit button (will be replaced with a reload button)
        document.getElementById("submit").style.display = "none";
        if (document.getElementById("czar").checked) {
            // send czar response to players
            var submissions = document.querySelectorAll('input[name=players]');
            var winner = document.querySelectorAll('input[name=players]:checked');
            if (winner.length == 1) {
                // only one winner
                let msg = {};
                msg.kind = "choice"
                msg.winner = winner[0].payload.player_id
                msg.payload = []
                for (var i = 0; i < submissions.length; i++) {
                    msg.payload.push(submissions[i].payload)
                }
                var data = JSON.stringify(msg)
                
                // now we should notify the server of the winner...
                fetch('/endround/', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json; charset=utf-8'
                    },
                    body: data
                })
                    .then(response => {
                        console.log('status:', response.status)
                        conn.send(data); // notify the other players
                        location.reload() // this reloads czar page
                    })
                    .then(data => console.log(data) );
            }
            return false;
        } else {
            // send player's selected cards as json to the czar
            var selectedCard = document.querySelectorAll('input[name=cards]:checked');
            if (selectedCard.length > 0) {
                let msg = {};
                msg.kind = "submission"
                msg.payload = {}
                msg.payload.player_id = document.getElementById("player_id").value;
                msg.payload.cards = {}
                for (var i = 0; i < selectedCard.length; i++) {
                    index = selectedCard[i].value;
                    msg.payload.cards[index] = document.getElementById("lbl_" + index).innerText;
                }
                conn.send(JSON.stringify(msg));
                logMessage("<b>Wait for Czar selection.</b>")
            }
            return false;
        }
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            logMessage("<b>Connection closed.</b>");
        };

        conn.onmessage = function (evt) {
            // even the submitting player receives this message
            czar = document.getElementById("czar").checked;
            
            const obj = JSON.parse(evt.data);
            switch(obj.kind) {
                case 'submission':
                    if(czar) {
                        // parse players selections
                        var chk = document.createElement("input");
                        chk.required = true;
                        chk.type = "radio";
                        chk.name = "players";
                        chk.id = "player_" + obj.payload.player_id;
                        chk.value = obj.payload.player_id;
                        chk.payload = obj.payload; // (ab)use custom DOM property
                        document.getElementById("white_cards").appendChild(chk)
                        for (let card_index in obj.payload.cards) {
                            var lbl = document.createElement("label")
                            lbl.id = "lbl_" + obj.payload.player_id + "_" + card_index
                            lbl.name = "lbl_"+ obj.payload.player_id
                            lbl.htmlFor = chk.id
                            lbl.innerText = obj.payload.cards[card_index]
                            document.getElementById("white_cards").appendChild(lbl)
                        }
                        // get number of players
                        fetch('/players/') // default GET
                            .then((response) => response.json()) // response.json() creates a variable and pipes to next then()
                            .then((json) => { // json is piped from previous then() and can be used here...
                                //console.log(json)
                                var submissions = document.querySelectorAll('input[name=players]');
                                document.getElementById("submit").disabled = submissions.length < (json.length - 1);
                            }); 
                    }
                    break;
                case 'choice':
                    if(!czar) {
                        //logMessage(evt.data)
                        // parse czar selections: print winner
                        winner = obj.winner
                        msg = "Winner: " // TODO: handle more winning messages
                        for (var i = 0; i < obj.payload.length; i++) {
                            if(obj.payload[i].player_id == winner) {
                                for (let card_index in obj.payload[i].cards) {
                                    msg += "\n"
                                    msg += obj.payload[i].cards[card_index]
                                }
                            }
                        }
                        logMessage(msg)
                        document.getElementById("next").removeAttribute("hidden")
                    }
                    break;
                default:
                    logMessage("ERROR")
                    break
            }      
        };

    } else {
        logMessage("<b>Your browser does not support WebSockets.</b>");
    }
};