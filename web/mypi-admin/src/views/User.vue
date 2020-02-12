<template>
    <span>
      <v-card>
        <v-card-test>Foo</v-card-test>
      </v-card>
    </span>
</template>

<script>
import axios from "axios";

export default {
    name: 'user',
    components: {
    },
    data() {
      return { user: {} }
    },
      mounted() {
        axios({ method: "GET", "url": "/api/users" }).then(result => {
                var users = result.data;
                for (var i = 0; i < users.length; i++) {
                  users[i].uri = "/users/"+users[i].text;
                }
                this.user = users;
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
