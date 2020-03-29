<template>
  <span>
    <v-container fluid>
      <v-stepper v-model="step" vertical>
        <v-stepper-step step="1" :complete="step > 1"
          >Alpine Version</v-stepper-step
        >
        <v-stepper-content step="1">
          Select the Alpine Linux version which shall be deployed on your
          Raspberry-PI
          <v-select
            :items="alpine_versions"
            label="Alpine Version"
            :value="alpine_version"
          ></v-select>
          Select the architecture. (aarch64 is recommended)
          <v-select
            :items="alpine_archs"
            label="Alpine Arch"
            :value="alpine_arch"
          ></v-select>
          <v-btn color="primary" @click="step = 2">Continue</v-btn>
        </v-stepper-content>
        <v-stepper-step step="2" :complete="step > 2"
          >System Preferences</v-stepper-step
        >
        <v-stepper-content step="2">
          <v-text-field
            name="hostname"
            label="Hostname"
            id="hostname"
            v-model="hostname"
          ></v-text-field>
          <!--
            <v-text-field prepend-icon="lock" name="password" label="root-Password" id="root_password" type="password" :value="root_password"></v-text-field>
            <v-container class="px-0" fluid>
              <v-checkbox label="Enable WiFi"></v-checkbox>
            </v-container>
            -->
          <v-btn color="primary" @click="step = 3">Continue</v-btn>

          <v-btn @click="step = 1">Back</v-btn>
        </v-stepper-content>
        <v-stepper-step step="3" :complete="step > 3">SD-Card</v-stepper-step>
        <v-stepper-content step="3">
          Select the SD-Card for your Alpine-Linux
          <v-select
            :items="sd_cards"
            label="SD-Card"
            :value="sd_card"
          ></v-select>
          <v-btn color="primary" :to="{name:'initializesd', params:{
            AlpineVersion: alpine_version,
            AlpineArch: alpine_arch,
            Hostname: hostname,
            Disk: sd_card
          }}">Continue</v-btn>
          <v-btn @click="step = 2">Back</v-btn>
        </v-stepper-content>
      </v-stepper>
    </v-container>
  </span>
</template>

<script>
import axios from "axios";

export default {
  name: "wizard-sd",
  components: {},
  data: function() {
    return {
      step: 1,
      alpine_versions: ["3.11.3", "3.11.2", "3.11.0"],
      alpine_version: "3.11.3",
      alpine_archs: ["aarch64", "armv7"],
      alpine_arch: "aarch64",
      hostname: "mypi",
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
  },
};
</script>
