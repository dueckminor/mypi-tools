<template>
  <v-container fluid>
    <div v-if="!haveData">
      <div class="text-center">
        <v-progress-circular indeterminate />
      </div>
    </div>
    <div v-else-if="hosts.length > 0">
      <v-card elevation="2" outlined v-for="(item, i) in hosts" :key="i">
        <v-card-title>{{item}}</v-card-title>
        <v-card-actions>
          <router-link :to="'/hosts/' + item + '/terminal'">
            <v-btn color="deep-purple lighten-2" text>
              Terminal
            </v-btn>
          </router-link>
          <router-link :to="'/hosts/' + item + '/wizardsd'">
            <v-btn color="deep-purple lighten-2" text>
              Create SD-Card
            </v-btn>
          </router-link>
          <router-link :to="'/hosts/' + item + '/actions/setup'">
            <v-btn color="deep-purple lighten-2" text>
              Connect to setup
            </v-btn>
          </router-link>
        </v-card-actions>
      </v-card>
      <v-btn dark absolute bottom right>
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </div>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import axios from "axios";

export default {
  name: "HostsView",
  components: {},
  computed: {},
  data: function() {
    return { haveData: false, hosts: [] };
  },
  mounted() {
    axios({ method: "GET", url: "/api/hosts" }).then((result) => {
      this.hosts = result.data;
      this.haveData = true;
    });
  },
  methods: {},
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
