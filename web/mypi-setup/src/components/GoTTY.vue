<template>
  <v-container fluid fill-height ref="xterm"/>
</template>

<script>
import { Terminal as Xterm } from 'xterm';
import { fit } from 'xterm/lib/addons/fit/fit';
import 'xterm/lib/xterm.css';
import { WebTTY, protocols } from './GoTTY/webtty'
import { ConnectionFactory } from "./GoTTY/websocket";
export default {
    name: "XTerm",
    mounted() {
        this._xterm = new Xterm({fontFamily:"Hack, monospace"});
        this._xterm.open(this.$refs['xterm']);

        var loc = window.location, new_uri;
        if (loc.protocol === "https:") {
            new_uri = "wss:";
        } else {
            new_uri = "ws:";
        }
        new_uri += "//" + loc.host;
        new_uri += "/ws";

        this._factory = new ConnectionFactory(new_uri, protocols);
        this._wt = new WebTTY(this, this._factory, '', '');
        this.closer = this._wt.open();
        window.addEventListener('resize', this.handleResize)

        fit(this._xterm);
    },
    beforeDestroy() {
        window.removeEventListener('resize', this.handleResize)
        this._xterm.dispose();
    },

    methods: {
        write(...args) {
            this._xterm.write(...args);
        },
        clear(/*...args*/) {
            this._xterm.clear();
        },
        blur() {
            this._xterm.blur();
        },
        focus() {
            this._xterm.focus();
        },
        handleResize (/*event*/) {
          fit(this._xterm);
        },
        info() {
          return { columns: this._xterm.cols, rows: this._xterm.rows };
        },
        output(data) {
          this._xterm.write(decodeURIComponent(escape(data)));
        },
        showMessage() {
          //console.log("showMessage")
        },
        removeMessage() {
          //console.log("removeMessage")
        },
        setWindowTitle(/*title*/) {
          //console.log("setWindowTitle: " + title)
        },
        setPreferences(/*value*/) {
          //console.log("setPreferences")
        },
        onInput(callback) {
          this._xterm.on("data", (data) => {
            callback(data);
          });
        },
        onResize(callback) {
          this._xterm.on("resize", (data) => {
            callback(data.cols, data.rows);
          });
        },
        reset() {
          this.removeMessage();
          this._xterm.clear();
        },
        deactivate() {
          this._xterm.off("data");
          this._xterm.off("resize");
          this._xterm.blur();
        },
        close() {
          this._xterm.destroy();
        }

    }
}
</script>

<style>
</style>