<template>
  <v-dialog v-model="show" max-width="500px" persistent>
    <!-- <v-btn text class="success">Login</v-btn> -->
<!--    <v-flex xs12 sm8 md4> -->
      <v-card class="elevation-12">
        <v-toolbar dark color="primary">
        <v-toolbar-title>Login</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-tooltip bottom>
        </v-tooltip>
        </v-toolbar>
        <v-card-text>
        <v-form>
            <v-text-field v-model="username" prepend-icon="person" name="login" label="Login" type="text" v-on:keyup.enter="$refs.passwordField.focus()"></v-text-field>
            <v-text-field v-model="password" prepend-icon="lock" name="password" label="Password" id="password" type="password" ref="passwordField" v-on:keyup.enter="loginClicked"></v-text-field>
        </v-form>
        </v-card-text>
        <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="loginClicked">Login</v-btn>
        </v-card-actions>
      </v-card>
 <!--   </v-flex> -->
  </v-dialog>
</template>

<script>
export default {
  props: {
     value: Boolean,
  },
  data() {
    return {
        username: "",
        password: ""
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
    loginClicked: function() {
      this.$emit('login-clicked',{username: this.username, password: this.password})
    },
    activated: function() {
      this.username = ""
      this.password = ""
    }
  }
}
</script>