<template>
  <v-app light id="app">
    <v-toolbar app></v-toolbar>
    <v-content>
      <v-container fluid>
        <v-card class="elevation-12">
          <v-card-text>
          </v-card-text>
          <v-btn color="accent" large @click.stop="showLogin=true">
            foo
          </v-btn>
        </v-card>
      </v-container>
    </v-content>
    <v-footer app></v-footer>
    <Login v-model="showLogin"/>
  </v-app>
</template>

<script>
//import GoTTY from './components/GoTTY'
import Login from './Login'

export default {
  name: 'app',
  components: {
    //GoTTY
    Login
  },
  data: function () {
    return {
      user: 'jochen',
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
