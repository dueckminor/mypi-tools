"use strict";

export const protocols = ["webtty"];
export const msgInputUnknown = '0';
export const msgInput = '1';
export const msgPing = '2';
export const msgResizeTerminal = '3';
export const msgUnknownOutput = '0';
export const msgOutput = '1';
export const msgPong = '2';
export const msgSetWindowTitle = '3';
export const msgSetPreferences = '4';
export const msgSetReconnect = '5';

export class WebTTY {
    constructor(term, connectionFactory, args, authToken) {
        this.term = term;
        this.connectionFactory = connectionFactory;
        this.args = args;
        this.authToken = authToken;
        this.reconnect = -1;
    }
    open() {
        let _this = this;
        let connection = this.connectionFactory.create();
        let pingTimer;
        let reconnectTimeout;
        let setup = function () {
            connection.onOpen(function () {
                let termInfo = _this.term.info();
                connection.send(JSON.stringify({
                    Arguments: _this.args,
                    AuthToken: _this.authToken
                }));
                let resizeHandler = function (colmuns, rows) {
                    connection.send(msgResizeTerminal + JSON.stringify({
                        columns: colmuns,
                        rows: rows
                    }));
                };
                _this.term.onResize(resizeHandler);
                resizeHandler(termInfo.columns, termInfo.rows);
                _this.term.onInput(function (input) {
                    connection.send(msgInput + input);
                });
                pingTimer = setInterval(function () {
                    connection.send(msgPing);
                }, 30 * 1000);
            });
            connection.onReceive(function (data) {
                let payload = data.slice(1);
                switch (data[0]) {
                    case msgOutput:
                        _this.term.output(atob(payload));
                        break;
                    case msgPong:
                        break;
                    case msgSetWindowTitle:
                        _this.term.setWindowTitle(payload);
                        break;
                    case msgSetPreferences:
                        _this.term.setPreferences(JSON.parse(payload));
                        break;
                    case msgSetReconnect:
                        //console.log("Enabling reconnect: " + autoReconnect + " seconds")
                        _this.reconnect = JSON.parse(payload);
                        break;
                }
            });
            connection.onClose(function () {
                clearInterval(pingTimer);
                _this.term.deactivate();
                _this.term.showMessage("Connection Closed", 0);
                if (_this.reconnect > 0) {
                    reconnectTimeout = setTimeout(function () {
                        connection = _this.connectionFactory.create();
                        _this.term.reset();
                        setup();
                    }, _this.reconnect * 1000);
                }
            });
            connection.open();
        };
        setup();
        return function () {
            clearTimeout(reconnectTimeout);
            connection.close();
        };
    }
}
