<template>
    <span>
    <v-list>
      <!--<v-subheader>Users</v-subheader>-->
      <v-list-item-group v-model="users" color="primary">
        <v-list-item
          v-for="(user, i) in users"
          :key="i"
          :to="user.uri"
        >
          <v-list-item-icon>
            <v-icon v-text="user.icon"></v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="user.text"></v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
    </span>
</template>

<script>
import axios from "axios";

export default {
    name: 'users',
    components: {
    },
    data() {
      return { users: [] }
    },
      mounted() {
        axios({ method: "GET", "url": "/api/users" }).then(result => {
                var users = result.data;
                for (var i = 0; i < users.length; i++) {
                  users[i].uri = "/users/"+users[i].text;
                }
                this.users = users;
            }, error => {
              if (error != null) {
                error = null
              }
              this.users = [
        {icon:"mdi-account-star",text:"Admin"},
        {icon:"mdi-account",text:"Jochen"}
        ] 
              //  console.error(error);
            });
    }
};
</script>
