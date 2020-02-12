<template>
  <v-app light id="app">
    <v-app-bar app></v-app-bar>
    <v-content>
      <v-container fluid>
        <!--
        <v-card class="elevation-12">
          <v-card-text>
            <GoTTY/>
          </v-card-text>
        </v-card>
        -->
        <v-stepper v-model="step" vertical>
          <v-stepper-step step="1" :complete="step > 1">Alpine Version</v-stepper-step>
          <v-stepper-content step="1">
              Select the Alpine Linux version which shall be deployed on your Raspberry-PI
              <v-select :items="alpine_versions" label="Alpine Version" :value="alpine_version"></v-select>
              Select the architecture. (aarch64 is recommended)
              <v-select :items="alpine_archs" label="Alpine Arch" :value="alpine_arch"></v-select>
              <v-btn color="primary" @click="step=2">Continue</v-btn>
          </v-stepper-content>
          <v-stepper-step step="2" :complete="step > 2">System Preferences</v-stepper-step>
          <v-stepper-content step="2">
            <v-text-field name="hostname" label="Hostname" id="root_password" :value="hostname"></v-text-field>
            <v-text-field prepend-icon="lock" name="password" label="root-Password" id="root_password" type="password" :value="root_password"></v-text-field>
            <v-container class="px-0" fluid>
              <v-checkbox label="Enable WiFi"></v-checkbox>
            </v-container>
            <v-btn color="primary" @click="step=3">Continue</v-btn>


            <v-btn @click="step=1">Back</v-btn>
          </v-stepper-content>
          <v-stepper-step step="3" :complete="step > 3">Features</v-stepper-step>
          <v-stepper-content step="3">
            <v-btn color="primary" @click="step=4">Continue</v-btn>
            <v-btn @click="step=2">Back</v-btn>
          </v-stepper-content>
          <v-stepper-step step="4" :complete="step > 4">SD-Card</v-stepper-step>
          <v-stepper-content step="4">
            <v-btn color="primary" @click="step=5">Continue</v-btn>
            <v-btn @click="step=3">Back</v-btn>
          </v-stepper-content>
        </v-stepper>
      </v-container>
    </v-content>
    <v-footer app></v-footer>
  </v-app>
</template>

<script>
//import GoTTY from './components/GoTTY'

export default {
  name: 'app',
  components: {
    //GoTTY
  },
  data: function () {
    return {
      step: 1,
      alpine_versions: ["3.11.3","3.11.2","3.11.0"],
      alpine_version: "3.11.3",
      alpine_archs: ["aarch64","armv7"],
      alpine_arch: "aarch64",
      hostname: "mypi",
      root_password: "",
      user: 'jochen',
      wifi: true,
      password: '',
      passcode: '',
      showLogin: true
    }
  },
  methods: {
    focusNext: function (next) {
      alert(JSON.stringify(this))
      alert(next)
    },
    setFocus: function (id) {
      document.getElementById(id).focus()
    },
    submit: function () {
      this.$refs.form.validate()
      fetch("/api/login",{
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          user: this.user,
          password: this.password,
          passcode: this.passcode
        })
      })
      alert("submit")
    }
  }
}
</script>

<style>
</style>
