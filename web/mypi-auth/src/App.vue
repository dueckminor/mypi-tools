<template>
  <v-app light id="app">
    <v-app-bar app>
      <v-btn class="success" @click.stop="showLoginDialog" :disabled="username!==''">
        Login
      </v-btn>
      <v-btn class="error" @click.stop="logout()" :disabled="username===''">
        Logout
      </v-btn>
      <span>{{username}}</span>
    </v-app-bar>
    <LoginDialog ref="loginDialog" v-model="showLogin" v-on:login-clicked="login"/>
  </v-app>
</template>

<script>
import LoginDialog from './LoginDialog'
//import VueRouter from 'vue-router'

export default {
  name: 'app',
  components: {
    LoginDialog
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
          console.log('login succeeded')
          if (response.status !== 200) {
            //noop
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
    const url = new URL(window.location.href)
    this.redirectURI = url.searchParams.get('redirect_uri')
    this.clientId = url.searchParams.get('params.client_id')
    this.responseType = url.searchParams.get('response_type')
    if (this.redirectURI != undefined && this.redirectURI !== "") {
      this.$router.replace("/")
      window.redirectURI = this.redirectURI
    }
    fetch("/status").then(res => res.json())
    .then(response => {
      if (response.username === "") {
        this.showLogin = true
      } else {
        this.username = response.username
      }
    })
  }
}
</script>

<style>
</style>
