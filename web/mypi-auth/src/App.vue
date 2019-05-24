<template>
  <v-app light id="app">
    <v-toolbar app>
      <v-btn class="success" @click.stop="showLoginDialog" :disabled="username!==''">
        Login
      </v-btn>
      <v-btn class="error" @click.stop="logout()" :disabled="username===''">
        Logout
      </v-btn>
      <v-text>{{username}}</v-text>
    </v-toolbar>
    <Login ref="loginDialog" v-model="showLogin" v-on:login-clicked="login"/>
  </v-app>
</template>

<script>
import Login from './Login'
import VueRouter from 'vue-router'

export default {
  name: 'app',
  components: {
    Login
  },
  data: function () {
    return {
      username: '',
      showLogin: false,
      clientId: '',
      responseType: '',
      redirectURI: '',
    }
  },
  methods: {
    logout: function() {
      fetch("/logout",{
        method: "POST"
      }).then((response) => {
        if (response.status === 202) {
          this.username = ""
        }
      })
    },
    login: function(params) {
      fetch("/login",{
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(params)
      }).then((response) => {
          if (response.status !== 200) {

          } else {
            if (window.redirectURI != undefined && window.redirectURI !== "") {
              window.location = "oauth/authorize?client_id="+encodeURIComponent(this.clientId) +
                "&redirect_uri=" + encodeURIComponent(this.redirectURI) +
                "&response_type=" + encodeURIComponent(this.response_type)
            } else {
              this.username = params.username
              this.showLogin = false
            }
          }
        }
      )
    },
    showLoginDialog: function() {
      this.$refs.loginDialog.password = ""
      this.showLogin = true
    }
  },
  beforeMount() {
    this.redirectURI = this.$route.query.redirect_uri
    this.clientId = this.$route.query.client_id
    this.responseType = this.$route.query.response_type
    if (this.redirectURI !== "") {
      this.$router.replace("/")
      window.redirectURI = this.redirectURI
    }
    fetch("/status").then(res => res.json())
    .then(response => {
      if (response.Username === "") {
        this.showLogin = true
      } else {
        this.username = response.Username
      }
    })
  }
}
</script>

<style>
</style>
