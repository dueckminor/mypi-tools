<template>
  <v-app theme="dark">
    <v-main>
      <v-row fill-width>
        <v-col/>
        <v-col>
          <v-container fluid fill-height fill-width align-center>
            <v-layout align-center justify-center>
              <v-card v-if="authenticated" class="elevation-12" max-height="500px" align-center style="margin: 25px;">
                <v-toolbar dark color="primary">
                  <v-toolbar-title>Logged in as '{{username}}'</v-toolbar-title>
                </v-toolbar>
                <div style="width:450px"></div>
                <v-card-actions>
                  <v-btn color="primary" @click="logout">Logout</v-btn>
                </v-card-actions>
              </v-card>
              <v-card v-if="!authenticated" class="elevation-12" max-height="500px" align-center style="margin: 20px;">
                <v-toolbar dark color="primary">
                  <v-toolbar-title>Login</v-toolbar-title>
                </v-toolbar>
                <v-card-text>
                  <v-form>
                    <div style="width:450px"></div>
                    <v-text-field v-model="username" prepend-icon="mdi-account" name="login" label="Login" type="text" v-on:keyup.enter="$refs.passwordField.focus()"></v-text-field>
                    <v-text-field v-model="password" prepend-icon="mdi-lock" name="password" label="Password" id="password" type="password" ref="passwordField" v-on:keyup.enter="login"></v-text-field>
                  </v-form>
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="primary" @click="login">Login</v-btn>
                  <v-spacer></v-spacer>
                </v-card-actions>
              </v-card>
              <v-spacer/>
            </v-layout>
          </v-container>
        </v-col>
        <v-col/>
      </v-row>
      <v-snackbar v-model="snackbar">
        {{ snackbar_text }}
      </v-snackbar>
    </v-main>
  </v-app>
</template>

<script>
export default {
  props: {
     value: Boolean,
  },
  data() {
    return {
        username: "",
        password: "",
        authenticated: false,
        clientId: '',
        responseType: '',
        redirectURI: '',
        snackbar: false,
        snackbar_text: '',
    }
  },
  computed: {
    show: {
      get () {
        return this.value
      },
      set (value) {
         this.$emit('input', value)
      }
    }
  },
  methods: {
    logout: function() {
      fetch("/logout",{
        method: "POST"
      }).then((response) => {
        if (response.status === 202) {
          this.authenticated = false
        }
      })
    },
    login: function() {
      fetch("/login",{
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({username: this.username, password: this.password})
      }).then((response) => {
          console.log('login succeeded')
          if (response.status !== 200) {
            this.snackbar_text = 'login failed'
            this.snackbar = true
          } else {
            if (window.redirectURI != undefined && window.redirectURI !== "") {
              window.location = "oauth/authorize?client_id="+encodeURIComponent(this.clientId) +
                "&redirect_uri=" + encodeURIComponent(this.redirectURI) +
                "&response_type=" + encodeURIComponent(this.response_type)
            } else {
              this.authenticated = true
              this.password = ""
            }
          }
        }
      )
    },
    activated: function() {
      this.username = ""
      this.password = ""
    },
  },
  beforeMount() {
    console.log('beforeMount')
    const url = new URL(window.location.href)
    this.redirectURI = url.searchParams.get('redirect_uri')
    console.log(this.redirectURI)
    this.clientId = url.searchParams.get('params.client_id')
    console.log(this.clientId)
    this.responseType = url.searchParams.get('response_type')
    console.log(this.responseType)
    if (this.redirectURI != undefined && this.redirectURI !== "") {
      this.$router.replace("/")
      window.redirectURI = this.redirectURI
    }
    fetch("/status").then(res => res.json())
    .then(response => {
      if (response.username === "") {
        this.authenticated = false
      } else {
        this.username = response.username
        this.authenticated = true
      }
    })
  }
}
</script>