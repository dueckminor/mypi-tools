<template>
  <v-dialog v-model="show" max-width="500px">
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
            <v-text-field v-model="username" prepend-icon="person" name="login" label="Login" type="text"></v-text-field>
            <v-text-field v-model="password" prepend-icon="lock" name="password" label="Password" id="password" type="password"></v-text-field>
        </v-form>
        </v-card-text>
        <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click.stop="submit()">Login</v-btn>
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
        username: "jochen",
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
      submit: function () {
      //this.$refs.form.validate()
      fetch("/api/login",{
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          username: this.username,
          password: this.password
        })
      })
      alert("submit")
    }
  }
}
</script>