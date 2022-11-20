<template>
  <v-container fluid>
    <h3>Alpine Version</h3>
    Select the Alpine Linux version which shall be deployed on your
    Raspberry-PI

    <v-select v-model="alpine_version" :items="alpine_versions" label="Alpine Version" :value="alpine_version"
      @change='versionChanged()'></v-select>
    Select the architecture. (aarch64 is recommended)
    <v-select v-model="alpine_arch" :items="alpine_archs" label="Alpine Arch" :value="alpine_arch"></v-select>
    
    <h3>Medium</h3>
    Select the SD-Card for your Alpine-Linux
    
    <v-select :items="sd_cards" label="SD-Card" :value="sd_card"></v-select>
  <v-btn color="primary" :to="{
        name: 'initializesd', params: {
          AlpineVersion: this.alpine_version,
          AlpineArch: this.alpine_arch,
          Hostname: this.hostname,
          Disk: this.sd_card
        }
      }">Continue</v-btn>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  name: "wizard-sd",
  components: {},
  data: function() {
    return {
      alpine_versions: [],
      alpine_version: "",
      alpine_archs: [],
      alpine_arch: "",
      root_password: "",
      user: "jochen",
      wifi: true,
      password: "",
      passcode: "",
      sd_cards: [],
      sd_card: "",
      showLogin: true
    };
  },
  mounted() {
    this.hostname = this.$route.params.host;
    axios({ method: "GET", url: "/api/hosts/localhost/disks" }).then(
      result => {
        var sd_cards = [];
        var selectionOK = false;
        for (var i = 0; i < result.data.length; i++) {
          sd_cards.push(result.data[i].name);
          selectionOK |= result.data[i].name == this.sd_card;
        }
        this.sd_cards = sd_cards;
        if (!selectionOK) {
          this.sd_card = result.data.length > 0 ? result.data[0].name : "";
        }
      },
      error => {
        if (error != null) {
          error = null;
        }
        this.sd_cards = [];
      }
    );
    axios({ method: "GET", url: "/api/downloads/alpine" }).then(
      result => {
        var alpine_versions = [];
        var alpine_archs = [];
        var selection_version_ok = false
        var selection_arch_ok = false
        for (var i = 0; i < result.data.length; i++) {
          var alpine_version = result.data[i].tags["alpine-version"];
          if (!alpine_versions.includes(alpine_version)) {
            alpine_versions.push(alpine_version);
            selection_version_ok |= alpine_version == this.alpine_version;
          } 
          var alpine_arch = result.data[i].tags["alpine-arch"];
          if (!alpine_archs.includes(alpine_arch)) {
            alpine_archs.push(alpine_arch);
            selection_arch_ok |= alpine_arch == this.alpine_arch;
          }
        }
        this.alpine_versions = alpine_versions;
        this.alpine_archs = alpine_archs;
        if (!selection_version_ok) {
          this.alpine_version = this.alpine_versions[this.alpine_versions.length-1];
        }
        if (!selection_arch_ok) {
          this.alpine_arch = this.alpine_archs[0];
        }
      },
      error => {
        if (error != null) {
          error = null;
        }
        this.alpine_versions = [];
        this.alpine_archs = [];
      }

    );
  },
  methods: {
        versionChanged() {
            alert("foo");
        },
    },
};
</script>
