"use strict";

export class ConnectionFactory {
    constructor(url, protocols) {
        this.url = url;
        this.protocols = protocols;
    }
    create() {
        return new Connection(this.url, this.protocols);
    }
}

export class Connection {
    constructor(url, protocols) {
        this.bare = new WebSocket(url, protocols);
    }
    open() {
        // noop
    }
    close() {
        this.bare.close();
    }
    send(data) {
        this.bare.send(data);
    }
    isOpen() {
        return (this.bare.readyState == WebSocket.CONNECTING ||
            this.bare.readyState == WebSocket.OPEN)
    }
    onOpen(callback) {
        this.bare.onopen = function ( /*event*/) {
            callback();
        };
    }
    onReceive(callback) {
        this.bare.onmessage = function (event) {
            callback(event.data);
        };
    }
    onClose(callback) {
        this.bare.onclose = function ( /*event*/) {
            callback();
        };
    }
}
