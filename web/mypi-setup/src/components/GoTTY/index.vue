<template>
  <v-container fluid>
    <div ref="xterm"/>
  </v-container>
</template>

<script>
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";
import "hack-font/build/web/hack.css";
import { WebTTY, protocols } from "./webtty";
import { ConnectionFactory } from "./websocket";

export default {
  name: "go-tty",
  props: ["path"],
  mounted() {
    // eslint-disable-next-line no-console
    console.log("mounted");

    var that = this;
    that.$nextTick(function() {
      if (that._xterm) {
        return
      }
      // eslint-disable-next-line no-console
      console.log("creating xterm");
      that._xterm = new Terminal({ fontFamily: "Hack, monospace", rows:10 });
      that._fitAddon = new FitAddon();
      that._xterm.open(that.$refs["xterm"]);

      var loc = window.location,
        new_uri;
      if (loc.protocol === "https:") {
        new_uri = "wss:";
      } else {
        new_uri = "ws:";
      }
      new_uri += "//" + loc.host;
      new_uri += that.$props.path;

      that._factory = new ConnectionFactory(new_uri, protocols);
      that._wt = new WebTTY(that, that._factory, "", "");
      that.closer = that._wt.open();
      window.addEventListener("resize", that.handleResize);

      that._xterm.loadAddon(that._fitAddon);
      that._fitAddon.fit();
    });
  },
  beforeUnmount() {
    // eslint-disable-next-line no-console
    console.log("beforeUnmount");
    window.removeEventListener("resize", this.handleResize);
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
    handleResize(/*event*/) {
      this._fitAddon.fit();
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
      // eslint-disable-next-line no-console
      console.log("setPreferences");
    },
    onInput(callback) {
      this._xterm.onData(function(data) {
        callback(data);
      });
    },
    onResize(callback) {
      this._xterm.onResize(data => {
        callback(data.cols, data.rows);
      });
    },
    reset() {
      this.removeMessage();
      this._xterm.clear();
    },
    deactivate() {
      this._xterm.blur();
    },
    close() {
      // eslint-disable-next-line no-console
      console.log("close");
      this._xterm.destroy();
    }
  }
};
</script>

<style></style>
